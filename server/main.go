package main

import (
	"log"
)

const port = ":3000" // Port for the GRPC server

func main() {
	log.Println("Starting server...")

	// Environment variables
	if err := loadEnvVars(); err != nil {
		log.Println(err)
	}
	verifyEnvVars()

	// Init Postgres pool
	pg := createPostgresPool()
	log.Println("Created the Postgres pool")

	// Init Redis client
	//rdb := createRedisClient()
	//log.Println("Created the Redis client")

	// Run GRPC server
	server := createServer(pg)
	server.run()
}
