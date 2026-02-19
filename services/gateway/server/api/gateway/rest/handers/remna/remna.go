package remna

import (
	remnaclient "gateway/server/api/gateway/grpc/remna_client"
	"net/http"
	"net/mail"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func GetUsersByEmail(log *zap.Logger, client *remnaclient.RemnaGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		_, err := mail.ParseAddress(email)
		if err != nil {
			http.Error(w, "invalid email", http.StatusBadRequest)
		}
		resp, err := client.GetUsersByEmail(email)
		if err != nil {
			http.Error(w, "failed to get users", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		render.JSON(w, r, resp)
	}
}
