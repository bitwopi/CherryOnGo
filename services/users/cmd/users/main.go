package main

import (
	"os"
	"users/config"
	"users/server/api/grpc"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic("failed to get root path")
	}
	cfg := config.MustLoad(rootPath + "/config/dev.yaml")
	server := grpc.NewServer(
		cfg.DSN,
		cfg.JWTSecret,
		cfg.Redis.Addr,
		cfg.Redis.DBNum,
		cfg.RefreshTokenTTL,
		cfg.AccessTokenTTL,
	)
	server.Start("localhost:" + cfg.GRPC.Port)
}
