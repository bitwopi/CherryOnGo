package suites

import (
	"context"
	"testing"
	"users/config"
	pb "users/server/api/grpc/user"
	jwtmanager "users/server/jwt_manager"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	UserClient pb.UserServiceClient
	JWTManager jwtmanager.JWTManager
}

var (
	grpcHost = "localhost:"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoad("/home/antisperma/Desktop/CherryOnGo/services/users/config/local.yaml")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.NewClient(
		grpcHost+cfg.GRPC.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	client := pb.NewUserServiceClient(cc)
	jwtManager, err := jwtmanager.NewJWTManager(cfg.JWTSecret)
	if err != nil {
		t.Fatalf("failed to create jwtManager")
	}
	return ctx, &Suite{
		Cfg:        cfg,
		UserClient: client,
		JWTManager: *jwtManager,
	}
}
