package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func createRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	redisAddr := fmt.Sprintf("%s:6379", host)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// test Redis connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("redis ping error:", err)
	}

	return rdb
}

func createPostgresPool() *pgxpool.Pool {
	host, user, pass := dbConfig()
	dbURL := fmt.Sprintf("host=%s user=%s password=%s database=%s sslmode=disable", host, user, pass, user)

	// connect the db pool
	dbPool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}

	// test the connection
	_, err = dbPool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		log.Fatalln("Database connection error:", err)
	}

	return dbPool
}

// Grab the db config from the environment.
func dbConfig() (host, user, pass string) {
	host = os.Getenv("DB_HOST")
	user = os.Getenv("DB_USER")
	pass = os.Getenv("DB_PASS")

	return
}

// UserEmailExists will check the database for a user with the given email.
func UserEmailExists(dbPool *pgxpool.Pool, email string) (bool, error) {
	err := dbPool.QueryRow(context.Background(), "SELECT email FROM users WHERE email = $1", email).Scan(&email)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// UserNameExists will check the database for a user with the given name.
func UserNameExists(dbPool *pgxpool.Pool, name string) (bool, error) {
	err := dbPool.QueryRow(context.Background(), "SELECT name FROM users WHERE name = $1", name).Scan(&name)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
