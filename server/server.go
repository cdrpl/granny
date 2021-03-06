package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/asaskevich/govalidator"
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
	pg   *pgxpool.Pool
	rdb  *redis.Client
	room Room
	proto.UnimplementedAuthServer
	proto.UnimplementedRoomServer
}

// Create new GRPC server.
func createServer(pg *pgxpool.Pool, rdb *redis.Client) *Server {
	return &Server{pg: pg, rdb: rdb, room: newRoom()}
}

// SignUp is used for new user registrations
func (s *Server) SignUp(ctx context.Context, in *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	// Input validation
	err := validateSignUpRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Name must be unique
	nameExists, err := userNameExists(in.Name, s.pg)
	if err != nil {
		return nil, status.Error(codes.Internal, "query name error")
	} else if nameExists {
		return nil, status.Error(codes.AlreadyExists, "name already exists")
	}

	// Email must be unique
	emailExists, err := userEmailExists(in.Email, s.pg)
	if err != nil {
		return nil, status.Error(codes.Internal, "query email error")
	} else if emailExists {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash error")
	}

	// Create user struct
	user := createUser(in.Name, in.Email, string(hash))

	// Insert user
	if err := insertUser(user, s.pg); err != nil {
		return nil, status.Error(codes.Internal, "insert user error")
	}

	log.Printf("New user registration: {email:%v name:%v}\n", user.Email, user.Name)

	return &proto.SignUpResponse{}, nil
}

// SignIn allows users to sign in.
func (s *Server) SignIn(ctx context.Context, in *proto.SignInRequest) (*proto.SignInResponse, error) {
	// Input validation
	err := validateSignInRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := User{}

	// Fetch user data
	sql := "SELECT id, name, pass FROM users WHERE email = $1"
	err = s.pg.QueryRow(context.Background(), sql, in.Email).Scan(&user.ID, &user.Name, &user.Pass)
	if err == pgx.ErrNoRows {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "query error")
	}

	// Compare request pass to the hashed pass
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(in.Pass))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate auth token
	token, err := createAuthToken(user.ID, s.rdb)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate auth token error: %v", err)
	}

	log.Println("User has logged in", in.Email)

	// Create sign in response
	res := &proto.SignInResponse{
		Id:    int32(user.ID),
		Name:  user.Name,
		Token: token,
	}

	return res, nil
}

// GetRoom will return a map of users in the room.
func (s *Server) GetRoom(ctx context.Context, in *proto.GetRoomRequest) (*proto.GetRoomResponse, error) {
	res := &proto.GetRoomResponse{
		Users: make(map[int32]*proto.User),
	}

	for _, user := range s.room.users {
		res.Users[int32(user.id)] = &proto.User{Id: int32(user.id), Name: user.name}
	}

	return res, nil
}

// JoinRoom will allow a user to join a room.
func (s *Server) JoinRoom(ctx context.Context, in *proto.JoinRoomReq) (*proto.JoinRoomRes, error) {
	id, _, _ := extractUserIDAndToken(ctx)

	user, err := findUser(id, s.pg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "join room error: %v", err)
	}

	ru := newRoomUser(id, user.Name)

	err = s.room.joinRoom(ru)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "join room error: %v", err)
	}

	return &proto.JoinRoomRes{}, nil
}

// UserJoined streams a user whenever a user joins the room.
func (s *Server) UserJoined(req *proto.UserJoinedReq, stream proto.Room_UserJoinedServer) error {
	id, _, _ := extractUserIDAndToken(stream.Context())

	ru := s.room.getUser(id)

	for {
		select {
		case joined := <-ru.joined:
			stream.Send(&proto.User{Id: int32(joined.id), Name: joined.name})
			break

		case <-stream.Context().Done():
			log.Println("Stream ended")
			return nil
		}
	}
}

// Run the GRPC server.
func (s *Server) run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	uInterceptor := UnaryInterceptor{rdb: s.rdb}
	sInterceptor := StreamInterceptor{rdb: s.rdb}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(uInterceptor.auth),
		grpc.StreamInterceptor(sInterceptor.auth),
	)

	proto.RegisterAuthServer(grpcServer, s)
	proto.RegisterRoomServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// SignUpValidator is used to validate sign up requests.
type SignUpValidator struct {
	Name  string `valid:"required,maxstringlength(16)"`
	Email string `valid:"required,maxstringlength(255)"`
	Pass  string `valid:"required,minstringlength(8),maxstringlength(255)"`
}

// Sanitize and validate the sign up request.
func validateSignUpRequest(req *proto.SignUpRequest) (err error) {
	req.Name = govalidator.Trim(req.Name, "")
	req.Email = govalidator.Trim(req.Email, "")
	req.Name = strings.ToLower(req.Name)

	req.Email, err = govalidator.NormalizeEmail(req.Email)
	if err != nil {
		return
	}

	v := SignUpValidator{Name: req.Name, Email: req.Email, Pass: req.Pass}
	_, err = govalidator.ValidateStruct(v)
	return
}

// SignInValidator is used to validate sign in requests.
type SignInValidator struct {
	Email string `valid:"required,maxstringlength(255)"`
	Pass  string `valid:"required,minstringlength(8),maxstringlength(255)"`
}

// Sanitize and validate the sign in request.
func validateSignInRequest(req *proto.SignInRequest) (err error) {
	req.Email = govalidator.Trim(req.Email, "")

	req.Email, err = govalidator.NormalizeEmail(req.Email)
	if err != nil {
		return
	}

	v := SignInValidator{Email: req.Email, Pass: req.Pass}
	_, err = govalidator.ValidateStruct(v)
	return
}
