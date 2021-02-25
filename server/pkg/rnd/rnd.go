// Package rnd provides helper functions to generate random data.
package rnd

import (
	"crypto/rand"
	"fmt"
)

// GenerateToken generates a cryptographically secure base 16 token.
func GenerateToken(bytes int) (string, error) {
	buff := make([]byte, bytes) // buf will hold the randomly generated data

	// generate random bytes
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}

	// convert random bytes to base 16 string and return
	return fmt.Sprintf("%x", buff), nil
}

// GenerateRememberToken will generate a cryptographically secure base 16 token.
// The remember token has a series identifier and a token. Both of the tokens are 16 bytes long.
func GenerateRememberToken() (id, token string, err error) {
	id, err = GenerateToken(16)
	if err != nil {
		return
	}

	token, err = GenerateToken(16)
	if err != nil {
		return
	}

	return
}
