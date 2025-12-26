package grpc

import (
	"context"
	"log"
	"net"
	"time"
	pb "users/server/api/grpc/users"
	"users/server/db"
	jwtmanager "users/server/jwt_manager"
	redismanager "users/server/redis_manager"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pgManager    *db.PgManager
	jwtManager   *jwtmanager.JWTManager
	redisManager *redismanager.RedisManager
	sessionTTL   time.Duration
	jwtTTL       time.Duration
	pb.UnimplementedUserServiceServer
}

func NewServer(
	dsn string,
	jwtSecret string,
	redisAddr string,
	redisDB int,
	sessionTTL time.Duration,
	jwtTTL time.Duration) *Server {
	pg, err := db.NewManager(dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	jm, err := jwtmanager.NewJWTManager(jwtSecret)
	if err != nil {
		log.Fatalf("failed to create JWT manager: %v", err)
	}
	rm := redismanager.NewRedisManager(redisAddr, "", redisDB)
	return &Server{
		pgManager:    pg,
		jwtManager:   jm,
		redisManager: rm,
		sessionTTL:   sessionTTL,
		jwtTTL:       jwtTTL,
	}
}

func (r *Server) Start(bindURL string) error {
	log.Println("Starting gRPC server...")
	r.pgManager.Migrate()
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, r)
	lis, err := net.Listen("tcp", bindURL)
	if err != nil {
		panic(err)
	}
	log.Println("gRPC server listening on: ", bindURL)
	return server.Serve(lis)
}

func (s *Server) SignUpUser(ctx context.Context, req *pb.AuthRequest) (*pb.SignUpResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "empty login or password")
	}
	userUUID, err := s.pgManager.CreateUser(req.Login, req.Password, 0, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return &pb.SignUpResponse{
		Status:   "user created",
		UserUuid: userUUID,
	}, nil
}

func (s *Server) AuthUser(ctx context.Context, req *pb.AuthRequest) (*pb.JWTResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "empty login or password")
	}
	user, err := s.pgManager.GetUserByEmail(req.Login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err := s.pgManager.CheckPassword(req.Login, req.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}
	jwt, err := s.jwtManager.NewJWT(user.UUID, s.jwtTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create JWT: %v", err)
	}
	refreshToken := s.jwtManager.NewRefreshToken()
	if _, err := s.redisManager.CreateSession(jwt, refreshToken, user.UUID, s.sessionTTL); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}
	log.Println(jwt)
	return &pb.JWTResponse{
		AccessToken:  jwt,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Server) RefreshJWT(ctx context.Context, req *pb.RefreshRequest) (*pb.JWTResponse, error) {
	if req.RefreshToken == "" || len(req.RefreshToken) != 36 {
		return nil, status.Error(codes.InvalidArgument, "invalid refresh token")
	}
	session, err := s.redisManager.GetSession(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get session: %v", err)
	}
	if session == nil {
		return nil, status.Error(codes.NotFound, "session not found")
	}
	if session.ExpiresAt.Before(time.Now()) {
		return nil, status.Error(codes.Unauthenticated, "refresh token expired")
	}
	newJWT, err := s.jwtManager.NewJWT(session.UserID, s.jwtTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create new JWT: %v", err)
	}
	if err := s.redisManager.DeleteSession(req.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete old session: %v", err)
	}
	refreshToken := s.jwtManager.NewRefreshToken()
	if _, err := s.redisManager.CreateSession(newJWT, refreshToken, session.UserID, s.sessionTTL); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}
	return &pb.JWTResponse{
		AccessToken:  newJWT,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Server) GetUserData(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetUserData not implemented")
}
