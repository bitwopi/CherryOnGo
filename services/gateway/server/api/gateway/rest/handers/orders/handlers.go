package handlers

import (
	pb "gateway/server/api/gateway/grpc/gen/order"
	orderclient "gateway/server/api/gateway/grpc/order_client"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

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

		render.JSON(w, r, resp)
	}
}

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

func GetOrder(log *zap.Logger, client *orderclient.OrderGRPCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUUID := r.URL.Query().Get("order_uuid")
		if orderUUID == "" {
			http.Error(w, "invalid parameter ", http.StatusBadRequest)
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
