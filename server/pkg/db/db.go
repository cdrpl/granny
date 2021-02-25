package db

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

const migrationDir = "./migration"

// CreateRedisClient will create a Redis client and return it.
func CreateRedisClient() *redis.Client {
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

// CreatePostgresPool will create a Postgres pool and return it.
func CreatePostgresPool() *pgxpool.Pool {
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

// MigrateUp will run the up migrations.
func MigrateUp(pgPool *pgxpool.Pool) error {
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

		// Read the file
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		// Execute the SQL
		sql := string(bytes)
		_, err = pgPool.Exec(context.Background(), sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckAuth will return true if the user is authorized and false if not.
// Authorization is determined by checking redis for the user id and token.
// The user id is the key and the token is the value.
func CheckAuth(rdb *redis.Client, userID, token string) (bool, error) {
	result, err := rdb.Get(context.Background(), userID).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("CheckAuth() failed: %v", err)
	}
	return result == token, nil
}

// User models the "users" table
type User struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Pass      string    `json:"pass"`
	CreatedOn time.Time `json:"createdOn"`
}

// Insert will insert the user into the database
func (u *User) Insert(dbPool *pgxpool.Pool) error {
	tx, err := dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	// Insert users row
	var id int
	sql := "INSERT INTO users(name, email, pass, created_at) VALUES($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRow(context.Background(), sql, u.Name, u.Email, u.Pass, time.Now()).Scan(&id)
	if err != nil {
		return err
	}

	// Insert positions row
	_, err = tx.Exec(context.Background(), "INSERT INTO positions(id) VALUES($1)", id)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
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
