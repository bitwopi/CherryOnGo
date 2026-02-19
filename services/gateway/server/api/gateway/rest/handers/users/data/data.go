package data

import (
	userclient "gateway/server/api/gateway/grpc/user_client"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetUser(log *zap.Logger, client *userclient.UserGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUuid := chi.URLParam(r, "uuid")
		_, err := uuid.Parse(userUuid)
		if err != nil {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
		}
		resp, err := client.GetUser(userUuid)
		if err != nil {
			http.Error(w, "failed to get users", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		render.JSON(w, r, resp)
	}
}
