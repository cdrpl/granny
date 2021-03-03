package main

import (
	"flag"
	"log"
)

const port = ":3000" // Port for the GRPC server

func main() {
	log.Println("Starting server")

	// Env vars
	if flags() == false {
		if err := loadEnvVars(); err != nil {
			log.Println(err)
		}
	}
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

func flags() bool {
	e := flag.Bool("e", false, ".env file will not be loaded if flag is given")
	flag.Parse()
	return *e
}
