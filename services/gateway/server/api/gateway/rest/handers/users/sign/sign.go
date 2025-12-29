package sign

import (
	userclient "gateway/server/api/gateway/grpc/user_client"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	JWT string `json:"jwt"`
}

func Auth(log *zap.Logger, client *userclient.UserGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := client.AuthUser(req.Login, req.Password)
		if err != nil {
			http.Error(w, "authentication failed: "+err.Error(), http.StatusUnauthorized)
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

func SignUp(log *zap.Logger, client *userclient.UserGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := client.SignUpUser(req.Login, req.Password)
		if err != nil {
			http.Error(w, "user registration failed: "+err.Error(), http.StatusInternalServerError)
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
