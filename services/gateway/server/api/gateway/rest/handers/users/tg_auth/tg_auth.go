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
	JWT string
}

// @Summary Аутентификация пользователя через тг
// @Description Возвращает jwt token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body Request true "Данные тг"
// @Success 200 {object} Response
// @Router /api/auth/telegram [post]
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

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			//TODO: set to true in production
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		render.JSON(w, r, Response{
			JWT: resp.AccessToken,
		})

	}
}
