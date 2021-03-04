package main

import (
	"context"
	"log"
	"net"

	"github.com/cdrpl/granny/server/proto"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.Internal, "query name error")
	} else if nameExists {
		return nil, status.Error(codes.AlreadyExists, "name already exists")
	}

	// Email must be unique
	emailExists, err := userEmailExists(email, s.pg)
	if err != nil {
		return nil, status.Error(codes.Internal, "query email error")
	} else if emailExists {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash error")
	}

	// Create user struct
	user := createUser(name, email, string(hash))

	// Insert user
	if err := insertUser(user, s.pg); err != nil {
		return nil, status.Error(codes.Internal, "insert user error")
	}

	log.Printf("New user registration: {email:%v name:%v}\n", user.Email, user.Name)

	return &proto.SignUpResponse{}, nil
}

// SignIn allows users to sign in.
func (s *Server) SignIn(ctx context.Context, in *proto.SignInRequest) (*proto.SignInResponse, error) {
	email := in.GetEmail()
	pass := in.GetPass()

	user := User{}

	// Fetch user data
	sql := "SELECT id, name, pass FROM users WHERE email = $1"
	err := s.pg.QueryRow(context.Background(), sql, email).Scan(&user.ID, &user.Name, &user.Pass)
	if err == pgx.ErrNoRows {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "query error")
	}

	// Compare request pass to the hashed pass
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate auth token
	token, err := createAuthToken(user.ID, s.rdb)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate auth token error: %v", err)
	}

	log.Println("User has logged in", email)

	// Create sign in response
	res := &proto.SignInResponse{
		Id:    int32(user.ID),
		Name:  user.Name,
		Token: token,
	}

	return res, nil
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
