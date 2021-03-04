package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func createRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	redisAddr := fmt.Sprintf("%s:6379", host)
	rdb := redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		DialTimeout: time.Millisecond * 500,
	})

	// Test Redis connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("redis ping error:", err)
	}

	return rdb
}

func createPostgresPool() *pgxpool.Pool {
	host, user, pass := dbConfig()
	dbURL := fmt.Sprintf("host=%s user=%s password=%s database=%s sslmode=disable", host, user, pass, user)

	// Connect the db pool
	dbPool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}

	// Test the connection
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

// Build the database
func dbUp(pg *pgxpool.Pool) error {
	items, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		// Open the file
		filePath := fmt.Sprintf("%s/%s", migrationDir, item.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Read the file
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		// Execute the SQL
		sql := string(bytes)
		_, err = pg.Exec(context.Background(), sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func userNameExists(name string, pg *pgxpool.Pool) (bool, error) {
	var id int

	err := pg.QueryRow(context.Background(), "SELECT id FROM users WHERE name = $1", name).Scan(&id)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func userEmailExists(email string, pg *pgxpool.Pool) (bool, error) {
	var id int

	err := pg.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func insertUser(user User, pg *pgxpool.Pool) error {
	sql := "INSERT INTO users (name, email, pass, created_at) VALUES ($1, $2, $3, $4)"
	_, err := pg.Exec(context.Background(), sql, user.Name, user.Email, user.Pass, user.CreatedAt)
	return err
}
