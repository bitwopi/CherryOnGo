package remna

import (
	pb "gateway/server/api/gateway/grpc/gen/remna"
	remnaclient "gateway/server/api/gateway/grpc/remna_client"
	"net/http"
	"net/mail"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type MultipleUsersResponseDTO struct {
	pb.MultipleUsersResponse
}

// @Summary Обновление статуса заказа
// @Description Возвращает объект заказа
// @Tags remna
// @Accept json
// @Produce json
// @Param email path string true "email пользователя"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} MultipleUsersResponseDTO
// @Router /api/remna/users/by-email/{email} [get]
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
