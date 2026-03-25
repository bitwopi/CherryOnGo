package handlers

import (
	pb "gateway/server/api/gateway/grpc/gen/order"
	orderclient "gateway/server/api/gateway/grpc/order_client"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary Создание заказа
// @Description Возвращает объект заказа
// @Tags order
// @Accept json
// @Produce json
// @Param order body pb.CreateOrderRequest true "Данные заказа"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} pb.OrderResponse
// @Router /api/order/create [post]
func CreateOrder(log *zap.Logger, client *orderclient.OrderGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.CreateOrderRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			log.Error(err.Error())
			return
		}
		var shopCard *pb.ShopCard
		if req.ShopCard != nil {
			shopCard = &pb.ShopCard{
				Uuid:        req.ShopCard.Uuid,
				Name:        req.ShopCard.Name,
				Description: req.ShopCard.Description,
				Category:    req.ShopCard.Category,
				Visible:     req.ShopCard.Visible,
			}
		}
		resp, err := client.CreateOrder(req.CustomerUuid, req.Status, shopCard, req.Price)
		if err != nil {
			http.Error(w, "failed to create order", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
		r.Response.StatusCode = http.StatusCreated
		render.JSON(w, r, resp)
	}
}

// @Summary Обновление статуса заказа
// @Description Возвращает объект заказа
// @Tags order
// @Accept json
// @Produce json
// @Param order body pb.OrderStatusRequest true "Запрос на изменение статуса"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} pb.OrderResponse
// @Router /api/order/update/status [post]
func UpdateOrderStatus(log *zap.Logger, client *orderclient.OrderGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.OrderStatusRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			log.Error(err.Error())
			return
		}
		resp, err := client.UpdateOrderStatus(req.OrderUuid, req.Status)
		if err != nil {
			http.Error(w, "failed to update order", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		render.JSON(w, r, resp)
	}
}

// @Summary Обновление статуса заказа
// @Description Возвращает объект заказа
// @Tags order
// @Accept json
// @Produce json
// @Param order_uuid path string true "uuid заказа"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} pb.OrderResponse
// @Router /api/order/get/{order_uuid} [get]
func GetOrder(log *zap.Logger, client *orderclient.OrderGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUUID := chi.URLParam(r, "order_uuid")
		_, err := uuid.Parse(orderUUID)
		if err != nil {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}
		resp, err := client.GetOrder(orderUUID)
		if err != nil {
			http.Error(w, "failed to get order", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		render.JSON(w, r, resp)
	}
}
