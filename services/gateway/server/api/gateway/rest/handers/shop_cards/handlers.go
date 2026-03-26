package shopcards

import (
	pb "gateway/server/api/gateway/grpc/gen/shop_card"
	shopcardclient "gateway/server/api/gateway/grpc/shopcard_client"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Request struct {
	Uuid        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Price       float32                `json:"price"`
	Visible     bool                   `json:"visible"`
	CoverUrl    string                 `json:"coverUrl"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type ShopCardDTO struct {
	pb.ShopCardResponse
}

type MultipleResponseDTO struct {
	pb.MultipleResponse
}

// @Summary Создание карточки товара
// @Description Возвращает объект карточки товара
// @Tags shop_card
// @Accept json
// @Produce json
// @Param shop_card body Request true "Данные карточки товара"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} ShopCardDTO
// @Router /api/shop_card/create [post]
func CreateShopCard(log *zap.Logger, client *shopcardclient.ShopCardGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			log.Error(err.Error())
			return
		}
		resp, err := client.CreateShopCard(
			req.Name,
			req.Description,
			req.Category,
			req.Price,
			req.Visible,
			req.CoverUrl,
			req.Metadata)
		if err != nil {
			http.Error(w, "failed to create shop card", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
		r.Response.StatusCode = http.StatusCreated
		render.JSON(w, r, resp)
	}
}

// @Summary Обновление карточки товара
// @Description Возвращает объект карточки товара
// @Tags shop_card
// @Accept json
// @Produce json
// @Param shop_card body Request true "Данные карточки товара"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} ShopCardDTO
// @Router /api/shop_card/update [post]
func UpdateShopCard(log *zap.Logger, client *shopcardclient.ShopCardGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			log.Error(err.Error())
			return
		}
		resp, err := client.UpdateShopCard(
			req.Uuid,
			req.Name,
			req.Description,
			req.Category,
			req.Price,
			req.Visible,
			req.CoverUrl,
			req.Metadata)
		if err != nil {
			http.Error(w, "failed to update shop card", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
		r.Response.StatusCode = http.StatusOK
		render.JSON(w, r, resp)
	}
}

// @Summary Получение карточки товара
// @Description Возвращает объект карточки товара
// @Tags shop_card
// @Accept json
// @Produce json
// @Param uuid path string true "uuid карточки товара"
// @Success 200 {object} ShopCardDTO
// @Router /api/shop_card/get/{uuid} [get]
func GetShopCard(log *zap.Logger, client *shopcardclient.ShopCardGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardUUID := chi.URLParam(r, "uuid")
		_, err := uuid.Parse(cardUUID)
		if err != nil {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}
		resp, err := client.GetShopCard(cardUUID)
		if err != nil {
			http.Error(w, "failed to get shop card", http.StatusNotFound)
			log.Error(err.Error())
			return
		}
		r.Response.StatusCode = http.StatusOK
		render.JSON(w, r, resp)
	}
}

// @Summary Получение всех карточек товаров
// @Description Возвращает список всех карточек товаров
// @Tags shop_card
// @Produce json
// @Success 200 {object} MultipleResponseDTO
// @Router /api/shop_card/get [get]
func GetAllShopCards(log *zap.Logger, client *shopcardclient.ShopCardGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.GetAllShopCards()
		if err != nil {
			http.Error(w, "failed to get shop cards", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
		r.Response.StatusCode = http.StatusOK
		render.JSON(w, r, resp)
	}
}
