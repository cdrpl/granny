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

	err := interceptorAuth(ctx, u.rdb)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// StreamInterceptor for authentication.
type StreamInterceptor struct {
	rdb *redis.Client
}

func (s *StreamInterceptor) auth(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := interceptorAuth(stream.Context(), s.rdb)
	if err != nil {
		return err
	}

	return handler(srv, stream)
}

// Shared auth code between stream and unary interceptor.
func interceptorAuth(ctx context.Context, rdb *redis.Client) error {
	id, token, err := extractUserIDAndToken(ctx)
	if err != nil {
		return status.Error(codes.Internal, "unauthorized")
	}

	// Verify auth token
	isValid, err := checkAuth(rdb, id, token)
	if err != nil {
		log.Printf("unary interceptor auth error: %v\n", err)
		return status.Error(codes.Internal, "an error has occured")
	} else if !isValid {
		return status.Error(codes.Unauthenticated, "unauthorized")
	}

	return nil
}
