package main

import (
	"flag"
	"orders/config"
	"orders/server/api/grpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfgPath := flag.String("c", "/config/dev.yaml", "cfg path")
	flag.Parse()
	cfg := config.MustLoad(*cfgPath)
	server, err := grpc.NewServer(cfg.DSN)
	if err != nil {
		panic(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.Start(cfg.GRPC.Host + ":" + cfg.GRPC.Port); err != nil {
			panic(err)
		}
	}()
	<-done
	server.Stop()
}
