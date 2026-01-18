package main

import (
	"flag"
	"os"
	"shopcards/config"
	"shopcards/server/api/grpc"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic("failed to get root path")
	}
	cfgPath := flag.String("c", "/config/dev.yaml", "cfg path")
	flag.Parse()
	cfg := config.MustLoad(rootPath + *cfgPath)
	server := grpc.NewServer(cfg.DSN)
	server.Start(cfg.GRPC.Host + ":" + cfg.GRPC.Port)
}
