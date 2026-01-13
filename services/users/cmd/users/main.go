package main

import (
	"flag"
	"os"
	"users/config"
	"users/server/api/grpc"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic("failed to get root path")
	}
	cfgPath := flag.String("c", "/config/dev.yaml", "cfg path")
	flag.Parse()
	cfg := config.MustLoad(rootPath + *cfgPath)
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
