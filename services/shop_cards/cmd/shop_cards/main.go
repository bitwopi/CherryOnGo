package main

import (
	"os"
	"shopcards/config"
	"shopcards/server/api/grpc"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic("failed to get root path")
	}
	cfg := config.MustLoad(rootPath + "/config/local.yaml")
	server := grpc.NewServer(cfg.DSN)
	server.Start(":" + cfg.GRPC.Port)
}
