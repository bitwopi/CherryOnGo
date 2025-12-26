package main

import (
	"context"
	"gateway/config"
	"gateway/server/api/gateway/rest/app"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg := config.MustLoad(rootPath + "/config/local.yaml")
	logger := setupLogger(cfg.Env)
	defer logger.Sync()
	app := app.NewApp(
		*cfg,
		logger,
	)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := app.Start(cfg.REST.Addr); err != nil {
			logger.Fatal("failed to start REST API server", zap.Error(err))
		}
	}()
	<-done
	logger.Info("shutting down gateway service")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		logger.Fatal("failed to stop REST API server", zap.Error(err))
	}
	logger.Info("gateway service stopped")
}

func setupLogger(env string) *zap.Logger {
	switch env {
	case "production":
		logger, err := zap.NewProduction()
		if err != nil {
			panic("failed to create production logger: " + err.Error())
		}
		return logger
	case "development":
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic("failed to create development logger: " + err.Error())
		}
		return logger
	default:
		logger := zap.NewExample()
		return logger
	}
}
