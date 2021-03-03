package main

import (
	"log"
)

const port = ":3000" // Port for the GRPC server

func main() {
	log.Println("Starting server")

	// Load env variables
	if err := loadEnvVars(); err != nil {
		log.Println(err)
	}

	// Verify env vars are set
	verifyEnvVars()

	// Init Postgres pool
	pg := createPostgresPool()
	log.Println("Postgres connected")

	// Init Redis client
	rdb := createRedisClient()
	log.Println("Redis connected")

	// Run GRPC server
	log.Printf("0.0.0.0%v\n", port)
	server := createServer(pg, rdb)
	server.run()
}
