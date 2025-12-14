package main

import (
	"remnawave/config"
	gs "remnawave/server/api/grpc"
)

func main() {
	cfg := config.NewConfig()
	// api := rest.NewAPIServer(cfg)
	// api.Start()
	gserver := gs.NewRemnaGRPCServer(cfg)
	gserver.Start(cfg.GRPCBindUrl)
}
