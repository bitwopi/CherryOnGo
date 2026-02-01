package main

import (
	"flag"
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
	cfgPath := flag.String("c", "/config/dev.yaml", "cfg path")
	flag.Parse()
	cfg := config.MustLoad(rootPath + *cfgPath)
	gserver := gs.NewRemnaGRPCServer(cfg)
	gserver.Start(cfg.GRPC.Host + ":" + cfg.GRPC.Port)
}
