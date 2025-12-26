package app

import (
	"context"
	"gateway/config"
	userclient "gateway/server/api/gateway/grpc/user_client"
	"gateway/server/api/gateway/rest/handers/users/refresh"
	"gateway/server/api/gateway/rest/handers/users/sign"
	jwtcheck "gateway/server/api/gateway/rest/middleware/jwt_check"
	jwtmanager "gateway/server/jwt_manager"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type App struct {
	client *userclient.UserGRPCClient
	jm     *jwtmanager.JWTManager
	logger *zap.Logger
	router *chi.Mux
	server *http.Server
}

func NewApp(cfg config.Config, logger *zap.Logger) *App {
	client, err := userclient.NewUserGRPCClient(
		cfg.GRPC.Addr,
		cfg.GRPC.Timeout,
		cfg.GRPC.MaxRetries)
	if err != nil {
		logger.Fatal("failed to create user gRPC client", zap.Error(err))
	}

	jm, err := jwtmanager.NewJWTManager(cfg.JWTSecret)
	if err != nil {
		logger.Fatal("failed to create JWT manager", zap.Error(err))
	}

	router := chi.NewRouter()
	server := &http.Server{
		Addr:         cfg.REST.Addr,
		ReadTimeout:  cfg.REST.RequestTimeout,
		WriteTimeout: cfg.REST.RequestTimeout,
		IdleTimeout:  cfg.REST.IdleTimeout,
	}
	return &App{
		client: client,
		jm:     jm,
		logger: logger,
		router: router,
		server: server,
	}
}

func (a *App) Start(addr string) error {
	a.SetupMiddleware()
	a.SetupRouter()
	a.server.Handler = a.router
	a.logger.Info("Starting REST API server", zap.String("bind_url", addr))
	// Here you would typically start the HTTP server
	return a.server.ListenAndServe()
}

func (a *App) Close() {
	a.client.Close()
}

func (a *App) SetupMiddleware() {
	a.router.Use(jwtcheck.New(*a.jm))
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
}

func (a *App) SetupRouter() {
	a.router.Post("/login", sign.Auth(a.logger, a.client))
	a.router.Post("/register", sign.SignUp(a.logger, a.client))
	a.router.Post("/auth/refresh", refresh.RefreshJWT(a.logger, a.client))
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.client.Close(); err != nil {
		return err
	}
	return a.server.Shutdown(ctx)
}
