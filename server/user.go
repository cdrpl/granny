package main

// User account data.
type User struct {
	ID    int64 `json:"id"`
	Email string
	Name  string `json:"name"`
	Pass  string
}
