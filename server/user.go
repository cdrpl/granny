package main

import "time"

// User account data.
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string
	Pass      string
	CreatedAt time.Time
}

func createUser(name, email, pass string) User {
	return User{
		Name:      name,
		Email:     email,
		Pass:      pass,
		CreatedAt: time.Now(),
	}
}
