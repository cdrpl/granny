package main

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata in request")
	}

	ids := md.Get("user-id")
	if ids == nil || len(ids) == 0 {
		return nil, status.Error(codes.Unauthenticated, "user-id metadata required")
	}

	// id must be an int
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user-id is not a valid number")
	}

	tokens := md.Get("token")
	if tokens == nil || len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "token metadata required")
	}

	// Verify auth token
	isValid, err := checkAuth(u.rdb, id, tokens[0])
	if err != nil {
		log.Println("unary interceptor auth error:", err.Error())
		return nil, status.Error(codes.Internal, "an error has occured")
	} else if !isValid {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	return handler(ctx, req)
}
