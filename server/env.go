package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// loadEnvVars will attempt to open the .env file and set the env vars.
func loadEnvVars() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse env file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "=")
		if len(split) == 2 && os.Getenv(split[0]) == "" {
			os.Setenv(split[0], split[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// verifyEnvVars will set env vars that haven't been set to default values.
func verifyEnvVars() {
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "development")
		log.Println("The ENV environment variable was not set, defaulting to development")
	}
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "127.0.0.1")
		log.Println("The DB_HOST environment variable was not set, defaulting to 127.0.0.1")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
		log.Println("The DB_USER environment variable was not set, defaulting to postgres")
	}
	if os.Getenv("DB_PASS") == "" {
		os.Setenv("DB_PASS", "password")
		log.Println("The DB_PASS environment variable was not set, defaulting to password")
	}
	if os.Getenv("REDIS_HOST") == "" {
		os.Setenv("REDIS_HOST", "127.0.0.1")
		log.Println("The REDIS_HOST environment variable was not set, defaulting to 127.0.0.1")
	}
}
