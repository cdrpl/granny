package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryInterceptor for authentication.
type UnaryInterceptor struct {
	rdb *redis.Client
}

// Used to authenticate requests.
func (u *UnaryInterceptor) auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/proto.Auth/SignIn" || info.FullMethod == "/proto.Auth/SignUp" {
		return handler(ctx, req)
	}

	id, token, err := extractUserIDAndToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "unauthorized")
	}

	// Verify auth token
	isValid, err := checkAuth(u.rdb, id, token)
	if err != nil {
		log.Printf("unary interceptor auth error: %v\n", err)
		return nil, status.Error(codes.Internal, "an error has occured")
	} else if !isValid {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	return handler(ctx, req)
}
