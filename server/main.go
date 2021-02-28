package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

const addr = "0.0.0.0:3010" // The address of the WebSocket server

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
	pg := createPostgresPool()
	log.Println("Created the Postgres pool")

	// Init Redis client
	rdb := createRedisClient()
	log.Println("Created the Redis client")

	// WebSocket server
	server := CreateServer()
	go server.run()

	// HTTP server
	runHTTPServer(server, rdb, pg)
}

func runHTTPServer(server *Server, rdb *redis.Client, pg *pgxpool.Pool) {
	log.Println("Server address -", addr)

	// WebSocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Health check and 404 handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.String() {
			case "/":
				fmt.Fprint(w, "OK")
				break

			case "/room":
				js, err := json.Marshal(server.getRoom())
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					break
				}

				fmt.Fprint(w, string(js))
				break

			default:
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// WebSocket upgrade handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")

		userID, err := verifyLogin(auth, server, rdb)
		if err != nil {
			log.Printf("upgrade WebSocket error for user %d: %v\n", userID, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Upgrade the connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade WebSocket error for user %d: %v\n", userID, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Fetch user's data
		user, err := queryUserData(userID, pg)
		if err != nil {
			log.Printf("upgrade WebSocket error for user %d: %v\n", userID, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Create WebSocket client
		client := CreateClient(userID, conn)
		go client.WritePump()
		go client.ReadPump(server)

		// Register client and user with the server
		server.Register(client, user)

		log.Println("User", userID, "has logged in")
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("Run HTTP server error:", err)
	}
}

// Error can be returned if user is already logged in, auth details are invalid, or if there is an issue with the Redis server.
// User ID will be returned if no error is encountered.
func verifyLogin(auth string, server *Server, rdb *redis.Client) (int64, error) {
	// Parse authentication string
	id, token, err := parseAuthorization(auth)
	if err != nil {
		return id, errors.New("invalid authorization header")
	}

	// Verify the authentication details
	isAuthorized, err := checkAuth(rdb, id, token)
	if err != nil {
		return id, err
	} else if !isAuthorized {
		return id, errors.New("auth token is not valid")
	}

	// Check if user client is already connected
	isOnline := server.IsUserOnline(id)
	if isOnline {
		return id, errors.New("User is already connected")
	}

	return id, nil
}
