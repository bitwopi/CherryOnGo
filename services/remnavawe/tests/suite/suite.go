package suite

import (
	"context"
	"os"
	"remnawave/config"
	pb "remnawave/server/api/grpc/remna"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg         *config.Config
	RemnaClient pb.RemnaServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	rootPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get root path: %v", err)
	}
	cfg := config.MustLoad(rootPath + "/../config/local.yaml")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.NewClient(
		cfg.GRPC.Host+cfg.GRPC.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	client := pb.NewRemnaServiceClient(cc)
	return ctx, &Suite{
		Cfg:         cfg,
		RemnaClient: client,
	}
}
