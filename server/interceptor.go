package main

import (
	"context"

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

	id := md.Get("user-id")
	if id == nil || len(id) == 0 {
		return nil, status.Error(codes.Unauthenticated, "user-id metadata required")
	}

	token := md.Get("token")
	if token == nil || len(token) == 0 {
		return nil, status.Error(codes.Unauthenticated, "token metadata required")
	}

	return handler(ctx, req)
}
