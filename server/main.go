package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/cdrpl/idlemon/pkg/env"
	"github.com/cdrpl/idlemon/pkg/game"
	"github.com/cdrpl/idlemon/pkg/ws"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// Persist player data to the database every saveInterval.
	saveInterval = time.Second * 60 * 5

	// Send player positions every posInterval.
	posInterval = time.Millisecond * 50

	// The player movement speed.
	moveSpeed = 0.5
)

var pgPool *pgxpool.Pool
var rdb *redis.Client

var server *ws.Server
var playerManager *game.PlayerManager

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	log.Println("Starting server...")

	// Environment variables
	loadEnvVars()
	env.VerifyEnvVars()

	// Init Postgres pool
	pgPool = createPostgresPool()
	log.Println("Created the Postgres pool")

	// Run the migrations
	/*err := db.MigrateUp(pgPool)
	if err != nil {
		log.Fatalln("Failed the run database migrations:", err)
	}*/

	// Init Redis client
	rdb = createRedisClient()
	log.Println("Created the Redis client")

	// WebSocket server
	server = ws.CreateServer()

	go handleIncomingData()

	// Player manager
	playerManager = game.CreatePlayerManager()

	// Save player data loop
	go savePlayerData()

	// Update player positions loop
	go updatePositions()

	// HTTP port
	port := os.Getenv("PORT")
	log.Println("Starting HTTP server on port", port)

	// HTTP controller/router
	controller := Controller{PgPool: pgPool, Rdb: rdb}
	router := httpRouter(controller)

	// Run HTTP server
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Println("Run HTTP server error:", err)
	}
}

// Loads the env vars from .env or .env.defaults if not in production.
func loadEnvVars() {
	if os.Getenv("ENV") != "production" {
		filename := env.LoadEnvVars()

		// Log the name of the .env file loaded
		if filename == "" {
			log.Println("Could not open the .env or .env.defaults file")
		} else {
			log.Println("Loaded env vars from", filename)
		}
	}
}

func httpRouter(c Controller) (r *gin.Engine) {
	// Create router based on env
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	} else {
		r = gin.Default()
	}

	// Health route
	r.GET("/health", c.health)

	// WebSocket upgrade handler route
	r.GET("/ws", upgradeHandler)

	// Sign up route
	r.POST("/sign-up", c.signUpHandler)

	// Sign in route
	r.POST("/sign-in", c.signInHandler)

	return
}

// Handles upgrading of WebSocket requests.
func upgradeHandler(c *gin.Context) {
	userID := c.GetHeader("user-id")
	token := c.GetHeader("auth-token")

	// Reject requests with missing header values
	if userID == "" || token == "" {
		c.String(401, "Unauthorized")
		return
	}

	log.Println("User", userID, "is attempting to login")

	// Convert userID to uint64
	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		log.Println("Failed to parse the user id:", err)
		c.String(401, "Unauthorized")
		return
	}

	// Check if player is already connected
	isOnline := server.PlayerOnline(uint32(userIDInt))
	if isOnline {
		log.Printf("Dual login detected for player '%v', the connection will be denied\n", userIDInt)
		c.String(401, "Unauthorized")
		return
	}

	// Verify the authentication token
	isAuthorized, err := checkAuth(rdb, userID, token)
	if err != nil {
		log.Println("Check auth failed:", err)
		c.String(401, "Unauthorized")
		return
	} else if !isAuthorized {
		log.Println("Auth token is not valid:", token)
		c.String(401, "Unauthorized")
		return
	}

	// Create the player structure
	player := game.Player{ID: uint32(userIDInt)}

	// If player data is present, clear the destination, else query the player data and register the player
	if playerManager.HasPlayer(player.ID) {
		playerManager.ClearPlayerDestination(player.ID)
		player, _ = playerManager.GetPlayerCopy(player.ID)
	} else {
		// Get the user data from the database
		err = pgPool.QueryRow(context.Background(), "SELECT name FROM users WHERE id = $1", userID).Scan(&player.Name)
		if err != nil {
			log.Println("Failed to query the users table:", err)
			c.String(401, "Unauthorized")
			return
		}

		// grab position data
		err = pgPool.QueryRow(context.Background(), "SELECT x, y FROM positions where id = $1", userID).Scan(&player.Position.X, &player.Position.Y)
		if err != nil {
			log.Println("Failed to query the positions table:", err)
			c.String(401, "Unauthorized")
			return
		}

		player.Destination = player.Position // To prevent player from moving to a zeroed out destination
		playerManager.Register(&player)
	}

	// Upgrade the connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade the connection:", err)
		c.String(401, "Unauthorized")
		return
	}

	// Create WebSocket client
	client := ws.CreateClient(player.ID, conn)
	server.Register(client)
	go client.WritePump()
	go client.ReadPump(server)

	// Send the player data back to the client
	err = sendPlayerDataToClient(player)
	if err != nil {
		server.Unregister(client)
		log.Println(err)
		c.String(401, "Unauthorized")
		return
	}

	broadcastPlayerConnected(player)
}

func savePlayerData() {
	saveTicker := time.NewTicker(saveInterval) // persist player data with this ticker
	defer saveTicker.Stop()

	for {
		<-saveTicker.C
		playerManager.SavePlayerData(server, pgPool)
	}
}

func updatePositions() {
	posTicker := time.NewTicker(posInterval) // position updates
	defer posTicker.Stop()

	for {
		<-posTicker.C
		playerManager.UpdatePlayerPositions(server, moveSpeed)
	}
}

// Send the player data in json format to the client.
func sendPlayerDataToClient(player game.Player) error {
	js, err := json.Marshal(player)
	if err != nil {
		return fmt.Errorf("Failed to convert the player data to JSON: %v", err)
	}

	// send player data to the client
	message := &ws.Message{Channel: ws.PlayerData, Data: js}
	go server.BroadcastSingle(message, player.ID)

	return nil
}

// Broadcast the player id and position on the player connected channel.
func broadcastPlayerConnected(player game.Player) {
	bytes := player.Pos()
	message := &ws.Message{Channel: ws.PlayerConnected, Data: bytes}
	go server.BroadcastAll(message)
	go playerManager.SendPlayerPositions(server, player.ID)
}

func handleIncomingData() {
	for {
		message := <-server.Incoming

		switch message.Channel {
		case ws.Chat:
			parseChatMessage(message)

			// re-broadcast the message to every other client
			server.BroadcastAllExclude(message, message.PlayerID)

		case ws.Destination:
			destination, err := parseDestinationMessage(message)
			if err != nil {
				log.Println(err)
				return
			}

			playerManager.SetPlayerDestination(message.PlayerID, destination)

		default:
			fmt.Fprintf(os.Stderr, "Received a message on an unsupported channel: %v\n", message.Channel)
		}
	}
}

// Create a new message with the player name prepended.
func parseChatMessage(message *ws.Message) {
	player, _ := playerManager.GetPlayerCopy(message.PlayerID)

	// create array to hold the player name
	var name [16]byte
	copy(name[:], []byte(player.Name))

	// append the player name and message data
	buf := make([]byte, 0)
	buf = append(buf, name[:]...)
	buf = append(buf, message.Data...)

	message.Data = buf
}

func parseDestinationMessage(message *ws.Message) (game.Vector, error) {
	vector := game.Vector{}

	if len(message.Data) != 8 {
		return vector, errors.New("Received message on the Destination channel that was not 8 bytes")
	}

	xInt := binary.LittleEndian.Uint32(message.Data[0:4])
	yInt := binary.LittleEndian.Uint32(message.Data[4:])

	x32 := math.Float32frombits(xInt)
	y32 := math.Float32frombits(yInt)

	vector.X = float64(x32)
	vector.Y = float64(y32)

	return vector, nil
}
