package jwtcheck

import (
	jwtmanager "gateway/server/jwt_manager"
	"net/http"
)

func New(jm jwtmanager.JWTManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			jwtHeader := r.Header.Get("Authorization")
			if jwtHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}
			jwtToken := jwtHeader[len("Bearer "):]
			if _, err := jm.ParseJWT(jwtToken); err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
