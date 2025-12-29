package app

import (
	"context"
	"gateway/config"
	userclient "gateway/server/api/gateway/grpc/user_client"
	"gateway/server/api/gateway/rest/handers/users/refresh"
	"gateway/server/api/gateway/rest/handers/users/sign"
	tgauth "gateway/server/api/gateway/rest/handers/users/tg_auth"
	checktgauth "gateway/server/api/gateway/rest/middleware/check_tg_auth"
	jwtmanager "gateway/server/jwt_manager"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type App struct {
	client   *userclient.UserGRPCClient
	jm       *jwtmanager.JWTManager
	logger   *zap.Logger
	router   *chi.Mux
	server   *http.Server
	botToken string
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
		client:   client,
		jm:       jm,
		logger:   logger,
		router:   router,
		server:   server,
		botToken: cfg.BotToken,
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
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
}

func (a *App) SetupRouter() {
	a.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/d3d1k/CherryOnGo/services/gateway/server/api/gateway/rest/app/index.html")
	})
	a.router.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", sign.Auth(a.logger, a.client))
		r.Post("/register", sign.SignUp(a.logger, a.client))
		r.Post("/refresh", refresh.RefreshJWT(a.logger, a.client))
	})
	a.router.Route("/api/auth/telegram", func(r chi.Router) {
		r.Use(checktgauth.New(a.botToken))
		r.Post("/", tgauth.New(a.logger, a.client))
	})

}

func (a *App) Stop(ctx context.Context) error {
	if err := a.client.Close(); err != nil {
		return err
	}
	return a.server.Shutdown(ctx)
}
