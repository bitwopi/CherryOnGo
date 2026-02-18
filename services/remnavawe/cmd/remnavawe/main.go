package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"remnawave/config"
	gs "remnawave/server/api/grpc"
	"syscall"
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
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := gserver.Start(cfg.GRPC.Host + ":" + cfg.GRPC.Port); err != nil {
			panic(err)
		}
	}()
	<-done
	gserver.Stop()

}
