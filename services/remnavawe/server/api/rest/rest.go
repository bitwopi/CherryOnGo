package rest

import (
	"encoding/json"
	"net/http"

	"remnawave/client"
	"remnawave/config"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type APIServer struct {
	Config *config.Config
	Logger *zap.Logger
	router *chi.Mux
	api    *client.Client
}

type CreateSubscriptionReq struct {
	Username   string `json:"username"`
	TelegramID string `json:"tg_id"`
	Duration   int    `json:"duration_days"`
}

func NewAPIServer(config *config.Config) *APIServer {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	client := client.NewClient(config.RemnaAPIKey, config.RemnaURL)
	return &APIServer{
		Config: config,
		Logger: logger,
		router: chi.NewRouter(),
		api:    client,
	}
}

func (s *APIServer) Start() error {
	// Implementation to start the REST API server
	s.configureRouter()
	s.Logger.Info("Starting REST API server", zap.String("bind_url", s.Config.RESTBindUrl))
	// Here you would typically start the HTTP server
	return http.ListenAndServe(s.Config.RESTBindUrl, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Get("/hello", s.Hello())
	s.router.Get("/remna/ping", s.PingRemna())
	s.router.Get("/user/create", s.CreateRemnaUser())
}

func (s *APIServer) Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}
}

func (s *APIServer) CreateRemnaUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if ct := r.Header.Get("Content-Type"); ct == "" || !strings.HasPrefix(ct, "application/json") {
		// 	http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		// 	return
		// }
		// const maxBody = 1 << 20 // 1 MiB
		// dec := json.NewDecoder(io.LimitReader(r.Body, maxBody))
		// dec.DisallowUnknownFields()
		// var req CreateSubscriptionReq
		// if err := dec.Decode(&req); err != nil {
		// 	http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		// 	return
		// }
		_, err := s.api.CreateUser(client.Plans["3:30"], "", "", "")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "failed", "error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func (s *APIServer) PingRemna() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.api.Ping()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "msg": "pong"})
		}
	}
}
