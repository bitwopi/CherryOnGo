package app

import (
	"context"
	"gateway/config"
	orderclient "gateway/server/api/gateway/grpc/order_client"
	remnaclient "gateway/server/api/gateway/grpc/remna_client"
	userclient "gateway/server/api/gateway/grpc/user_client"
	handlers "gateway/server/api/gateway/rest/handers/orders"
	"gateway/server/api/gateway/rest/handers/remna"
	"gateway/server/api/gateway/rest/handers/users/data"
	"gateway/server/api/gateway/rest/handers/users/refresh"
	"gateway/server/api/gateway/rest/handers/users/sign"
	tgauth "gateway/server/api/gateway/rest/handers/users/tg_auth"
	checktgauth "gateway/server/api/gateway/rest/middleware/check_tg_auth"
	jwtcheck "gateway/server/api/gateway/rest/middleware/jwt_check"
	jwtmanager "gateway/server/jwt_manager"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type App struct {
	userClient  *userclient.UserGRPCClient
	orderClient *orderclient.OrderGRPCClient
	remnaClient *remnaclient.RemnaGRPCClient
	jm          *jwtmanager.JWTManager
	logger      *zap.Logger
	router      *chi.Mux
	server      *http.Server
	botToken    string
}

func NewApp(cfg config.Config, logger *zap.Logger) *App {
	uClient, err := userclient.NewUserGRPCClient(
		cfg.UserService.Addr,
		cfg.UserService.Timeout,
		cfg.UserService.MaxRetries)
	if err != nil {
		logger.Fatal("failed to create user gRPC client", zap.Error(err))
	}

	oClient, err := orderclient.NewOrderGRPCClient(
		cfg.OrderService.Addr,
		cfg.OrderService.Timeout,
		cfg.OrderService.MaxRetries)
	if err != nil {
		logger.Fatal("failed to create order gRPC client", zap.Error(err))
	}

	rClient, err := remnaclient.NewRemnaGRPCClient(
		cfg.RemnaService.Addr,
		cfg.RemnaService.Timeout,
		cfg.RemnaService.MaxRetries,
	)

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
		userClient:  uClient,
		orderClient: oClient,
		remnaClient: rClient,
		jm:          jm,
		logger:      logger,
		router:      router,
		server:      server,
		botToken:    cfg.BotToken,
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
	a.userClient.Close()
}

func (a *App) SetupMiddleware() {
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
}

func (a *App) SetupRouter() {
	a.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	a.router.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", sign.Auth(a.logger, a.userClient))
		r.Post("/register", sign.SignUp(a.logger, a.userClient))
		r.Post("/refresh", refresh.RefreshJWT(a.logger, a.userClient))
	})
	a.router.Route("/api/auth/telegram", func(r chi.Router) {
		r.Use(checktgauth.New(a.botToken))
		r.Post("/", tgauth.New(a.logger, a.userClient))
	})
	a.router.Route("/api/users", func(r chi.Router) {
		r.Use(jwtcheck.New(*a.jm))
		r.Post("/{uuid}", data.GetUser(a.logger, a.userClient))
	})
	a.router.Route("/api/order", func(r chi.Router) {
		r.Use(jwtcheck.New(*a.jm))
		r.Post("/create", handlers.CreateOrder(a.logger, a.orderClient))
		r.Post("/update/status", handlers.UpdateOrderStatus(a.logger, a.orderClient))
		r.Get("/get", handlers.GetOrder(a.logger, a.orderClient))
	})
	a.router.Route("/api/remna}", func(r chi.Router) {
		r.Use(jwtcheck.New(*a.jm))
		r.Post("/users/by-email/{email}", remna.GetUsersByEmail(a.logger, a.remnaClient))
	})
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.userClient.Close(); err != nil {
		return err
	}
	return a.server.Shutdown(ctx)
}
