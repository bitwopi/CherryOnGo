package main

import (
	"users/config"
	"users/server/api/grpc"
)

func main() {
	cfg := config.MustLoad("/home/antisperma/Desktop/CherryOnGo/services/users/config/local.yaml")
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
