package data

import (
	"fmt"
	pb "gateway/server/api/gateway/grpc/gen/users"
	userclient "gateway/server/api/gateway/grpc/user_client"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserResponseDTO struct {
	pb.UserResponse
}

// @Summary Получение информации о пользователе
// @Description Возвращает объект пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "uuid пользователя"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} UserResponseDTO
// @Router /api/users/{uuid} [get]
func GetUser(log *zap.Logger, client *userclient.UserGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUuid := chi.URLParam(r, "uuid")
		_, err := uuid.Parse(userUuid)
		if err != nil {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
		}
		resp, err := client.GetUser(userUuid)
		msg := fmt.Sprintf(
			"%v, %v, %v, %v, %v, %v, %v, %v, %v",
			resp.UserUuid,
			resp.Email,
			resp.FirstName,
			resp.PhotoUrl,
			resp.Active,
			resp.ReferralUuid,
			resp.Roles,
			resp.TgId,
			resp.Trial,
		)
		log.Debug(msg)
		if err != nil {
			http.Error(w, "failed to get users", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		render.JSON(w, r, resp)
	}
}
