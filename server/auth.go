package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

// Extracts the user ID and auth token from an authorization header.
func parseAuthorization(authorization string) (id int64, token string, err error) {
	split := strings.Split(authorization, ":")

	if len(split) != 2 {
		return 0, "", errors.New("Invalid authorization header")
	}

	id, err = strconv.ParseInt(split[0], 10, 4)
	token = split[1]

	return
}

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