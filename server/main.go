package main

import (
	"flag"
	"log"
	"time"
)

const (
	port         = ":3000"            // Port for the GRPC server
	migrationDir = "./db"             // Directory that holds the SQL files
	roomSize     = 5                  // Max users in a room
	tokenBytes   = 16                 // Num bytes in the auth token, num chars in the token will be tokenBytes * 2
	tokenExpire  = time.Hour * 24 * 7 // Time till auth tokens expire
)

func main() {
	log.Println("Starting server")

	// Env vars
	if !flags() {
		if err := loadEnvVars(); err != nil {
			log.Println(err)
		}
	}
	verifyEnvVars()

	// Init Postgres pool
	pg := createPostgresPool()
	log.Println("Postgres connected")

	// Construct tables
	if err := dbUp(pg); err != nil {
		log.Fatal("db migrations error: ", err)
	}

	// Init Redis client
	rdb := createRedisClient()
	log.Println("Redis connected")

	// Run GRPC server
	log.Printf("0.0.0.0%v\n", port)
	server := createServer(pg, rdb)
	server.run()
}

func flags() bool {
	e := flag.Bool("e", false, ".env file will not be loaded if flag is given")
	flag.Parse()
	return *e
}
