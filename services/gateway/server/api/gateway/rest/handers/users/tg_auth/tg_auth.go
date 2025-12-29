package tgauth

import (
	userclient "gateway/server/api/gateway/grpc/user_client"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Request struct {
	AuthDate  time.Time
	FirstName string
	Hash      string
	ID        int64
	PhotoURL  string
	Username  string
}

type Response struct {
	JWTToken string
}

func New(log *zap.Logger, client *userclient.UserGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(err.Error())
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		resp, err := client.TgOAuth(req.ID, req.FirstName, req.PhotoURL, req.Username)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to authorize user", http.StatusBadRequest)
			return
		}
		render.JSON(w, r, resp)

	}
}
