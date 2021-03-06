package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// checkAuth will return true if the user is authorized and false if not.
// Authorization is determined by checking redis for the user id and token.
// The user id is the key and the token is the value.
func checkAuth(rdb *redis.Client, userID int64, token string) (bool, error) {
	result, err := rdb.Get(context.Background(), fmt.Sprintf("%d", userID)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("CheckAuth() failed: %v", err)
	}
	return result == token, nil
}

// createAuthToken will generate a random token and store it in Redis using the given id as the key.
// Return value (token, err).
func createAuthToken(id int64, rdb *redis.Client) (string, error) {
	token, err := generateToken(tokenBytes)
	if err != nil {
		return token, errors.New("generateToken error: " + err.Error())
	}

	err = rdb.Set(context.Background(), fmt.Sprint(id), token, tokenExpire).Err()
	if err != nil {
		return token, errors.New("redis set ex error: " + err.Error())
	}

	return token, nil
}
