package main

import (
	"context"
	"log"
	"net"

	"github.com/cdrpl/granny/server/proto"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server handles GRPC requests.
type Server struct {
	pg  *pgxpool.Pool
	rdb *redis.Client
	proto.UnimplementedAuthServer
}

// Create new GRPC server.
func createServer(pg *pgxpool.Pool, rdb *redis.Client) *Server {
	return &Server{pg: pg, rdb: rdb}
}

// SignUp is used for new user registrations
func (s *Server) SignUp(ctx context.Context, in *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	name := in.GetName()
	email := in.GetEmail()
	pass := in.GetPass()

	// Name must be unique
	nameExists, err := userNameExists(name, s.pg)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "query name error")
	} else if nameExists {
		return nil, grpc.Errorf(codes.AlreadyExists, "name already exists")
	}

	// Email must be unique
	emailExists, err := userEmailExists(email, s.pg)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "query email error")
	} else if emailExists {
		return nil, grpc.Errorf(codes.AlreadyExists, "email already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "hash error")
	}

	// Create user struct
	user := createUser(name, email, string(hash))

	// Insert user
	if err := insertUser(user, s.pg); err != nil {
		return nil, grpc.Errorf(codes.Internal, "insert user error")
	}

	log.Printf("New user registration: {email:%v name:%v}\n", user.Email, user.Name)

	return &proto.SignUpResponse{}, nil
}

// Run the GRPC server.
func (s *Server) run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
