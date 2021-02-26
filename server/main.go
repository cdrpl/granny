package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	addr = "0.0.0.0:3010" // The address of the WebSocket server
)

var pgPool *pgxpool.Pool
var rdb *redis.Client

var server *Server
var playerManager *PlayerManager

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	log.Println("Starting server...")

	// Environment variables
	filename := loadEnvVars()
	if filename == "" {
		log.Println("Could not open the .env or .env.defaults file")
	} else {
		log.Println("Loaded env vars from", filename)
	}
	verifyEnvVars()

	// Init Postgres pool
	pgPool = createPostgresPool()
	log.Println("Created the Postgres pool")

	// Init Redis client
	rdb = createRedisClient()
	log.Println("Created the Redis client")

	// WebSocket server
	server = CreateServer()

	go handleIncomingData()

	// Player manager
	playerManager = CreatePlayerManager()

	// HTTP server
	runHTTPServer()
}

func runHTTPServer() {
	log.Println("Server address -", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.String() == "/" {
				fmt.Fprint(w, "OK")
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

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
	player := Player{ID: id}

	// If player data is present, clear the destination, else query the player data and register the player
	if !playerManager.HasPlayer(player.ID) {
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
	client := CreateClient(player.ID, conn)
	server.Register(client)
	go client.WritePump()
	go client.ReadPump(server)
}

func handleIncomingData() {
	for {
		data := <-server.Incoming
		fmt.Printf("Received message %s", data)
	}
}
