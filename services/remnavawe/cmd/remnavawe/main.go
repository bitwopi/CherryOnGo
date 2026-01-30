package main

import (
	"log"
	"os"
	"remnawave/config"
	gs "remnawave/server/api/grpc"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Print(rootPath)
	cfg := config.MustLoad(rootPath + "/config/local.yaml")
	gserver := gs.NewRemnaGRPCServer(cfg)
	gserver.Start(cfg.GRPC.Host + cfg.GRPC.Port)
}
