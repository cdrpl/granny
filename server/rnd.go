package main

import (
	"crypto/rand"
	"encoding/hex"
)

// generateToken generates a cryptographically secure base 16 token.
func generateToken(bytes int) (string, error) {
	buff := make([]byte, bytes) // buf will hold the randomly generated data

	// generate random bytes
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}

	// convert random bytes to base 16 string and return
	return hex.EncodeToString(buff), nil
}
