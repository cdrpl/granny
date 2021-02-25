package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

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
