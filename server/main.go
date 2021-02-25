package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/cdrpl/granny/server/pkg/env"
	"github.com/cdrpl/granny/server/pkg/game"
	"github.com/cdrpl/granny/server/pkg/ws"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	addr = "0.0.0.0:3010"

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

	// HTTP server
	runHTTPServer()
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

func runHTTPServer() {
	log.Println("Starting server -", addr)

	http.HandleFunc("/ws", upgradeHandler)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("Run HTTP server error:", err)
	}
}

// Handles upgrading of WebSocket requests.
func upgradeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s/ws", addr)

	auth := r.Header.Get("authorization")
	id, token, err := parseAuthorization(auth)
	if err != nil {
		log.Println("Invalid authorization header:", auth)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("User", id, "is attempting to login")

	// Check if player is already connected
	isOnline := server.PlayerOnline(id)
	if isOnline {
		log.Println("User", id, "login failed since user is already online")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify the authentication token
	isAuthorized, err := checkAuth(rdb, id, token)
	if err != nil {
		log.Println("Check auth failed:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !isAuthorized {
		log.Println("Auth token is not valid:", token)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create the player structure
	player := game.Player{ID: id}

	// If player data is present, clear the destination, else query the player data and register the player
	if playerManager.HasPlayer(player.ID) {
		player, _ = playerManager.GetPlayerCopy(player.ID)
	} else {
		err = pgPool.QueryRow(context.Background(), "SELECT name FROM users WHERE id = $1", id).Scan(&player.Name)
		if err != nil {
			log.Println("Failed to query the users table:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		playerManager.Register(&player)
	}

	// Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade the connection:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
	fmt.Println("broadcastPlayerConnected not implemented yet")
	//bytes := player.Pos()
	//message := &ws.Message{Channel: ws.PlayerConnected, Data: bytes}
	//go server.BroadcastAll(message)
	//go playerManager.SendPlayerPositions(server, player.ID)
}

func handleIncomingData() {
	for {
		message := <-server.Incoming

		switch message.Channel {
		case ws.Chat:
			parseChatMessage(message)

			// re-broadcast the message to every other client
			server.BroadcastAllExclude(message, message.PlayerID)

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
