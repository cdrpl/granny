package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/metadata"
)

// checkAuth will return true if the user is authorized and false if not.
// Authorization is determined by checking redis for the user id and token.
// The user id is the key and the token is the value.
func checkAuth(rdb *redis.Client, userID int, token string) (bool, error) {
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
func createAuthToken(id int, rdb *redis.Client) (string, error) {
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

// extract user ID and auth token from the gRPC context.
func extractUserIDAndToken(ctx context.Context) (id int, token string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return id, token, errors.New("extract user id error: no metadata")
	}

	ids := md.Get("user-id")
	if ids == nil || len(ids) == 0 {
		return id, token, errors.New("user-id metadata required")
	}

	// Convert string id to int
	id, err = strconv.Atoi(ids[0])
	if err != nil {
		return id, token, errors.New("user-id is not a valid number")
	}

	tokens := md.Get("token")
	if tokens == nil || len(tokens) == 0 {
		return id, token, errors.New("token metadata required")
	}

	token = tokens[0]
	return
}
