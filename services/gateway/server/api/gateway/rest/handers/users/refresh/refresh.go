package refresh

import (
	userclient "gateway/server/api/gateway/grpc/user_client"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Request struct {
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	JWT string `json:"jwt"`
}

func RefreshJWT(log *zap.Logger, client *userclient.UserGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rToken, err := r.Cookie("refresh_token")
		if http.ErrNoCookie == err {
			http.Error(w, "refresh token not found", http.StatusUnauthorized)
			return
		}
		req := Request{RefreshToken: rToken.Value}
		log.Debug(rToken.Value)
		resp, err := client.RefreshJWT(req.RefreshToken)
		if err != nil {
			http.Error(w, "failed to refresh JWT: "+err.Error(), http.StatusUnauthorized)
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
