package grpc

import (
	"context"
	"log"
	"net"
	pb "orders/server/api/grpc/gen/order"
	"orders/server/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	db         db.DBManager
	grpcServer *grpc.Server
}

func NewServer(dsn string) (*Server, error) {
	manager, err := db.NewManager(dsn)
	if err != nil {
		return nil, err
	}
	return &Server{db: manager}, nil
}

func (s *Server) Start(bindURL string) error {
	s.db.Migrate()
	log.Println("Starting gRPC server...")
	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, s)
	lis, err := net.Listen("tcp", bindURL)
	if err != nil {
		return err
	}
	s.grpcServer = server
	log.Println("gRPC server listening on: ", bindURL)
	return server.Serve(lis)

}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	log.Println(req)
	if !isValidStatus(req.Status) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order status: %s", req.Status)
	}
	if req.ShopCard == nil {
		return nil, status.Errorf(codes.InvalidArgument, "shop card is required")
	}
	shopCard := db.ShopCard{
		UUID:        req.ShopCard.Uuid,
		Name:        req.ShopCard.Name,
		Description: req.ShopCard.Description,
		Category:    req.ShopCard.Category,
		Price:       &req.ShopCard.Price,
		Visible:     req.ShopCard.Visible,
		CoverURL:    &req.ShopCard.CoverUrl,
	}
	if req.ShopCard.CreatedAt != nil {
		shopCard.CreatedAt = req.ShopCard.CreatedAt.AsTime()
	}

	order, err := s.db.CreateOrder(
		req.CustomerUuid,
		db.OrderStatus(req.Status),
		&shopCard,
		req.Price,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	var sCard pb.ShopCard
	if order.ShopCard != nil {
		sCard = pb.ShopCard{
			Uuid:        order.ShopCard.UUID,
			Name:        order.ShopCard.Name,
			CreatedAt:   timestamppb.New(order.ShopCard.CreatedAt),
			Description: order.ShopCard.Description,
			Category:    order.ShopCard.Category,
			Visible:     order.ShopCard.Visible,
		}
		if order.ShopCard.Price != nil {
			sCard.Price = *order.ShopCard.Price
		}
		if order.ShopCard.CoverURL != nil {
			sCard.CoverUrl = *order.ShopCard.CoverURL
		}
	}
	return &pb.OrderResponse{
		Uuid:         order.UUID,
		CustomerUuid: order.CustomerUUID,
		Status:       string(order.Status),
		ShopCard:     &sCard,
		Price:        order.Price,
	}, nil
}

func (s *Server) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatusRequest) (*pb.OrderResponse, error) {
	if !isValidStatus(req.Status) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order status: %s", req.Status)
	}
	order, err := s.db.UpdateOrderStatus(
		req.OrderUuid,
		db.OrderStatus(req.Status),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}
	return &pb.OrderResponse{
		Uuid:         order.UUID,
		CustomerUuid: order.CustomerUUID,
		Status:       string(order.Status),
		ShopCard: &pb.ShopCard{
			Uuid:        order.ShopCard.UUID,
			Name:        order.ShopCard.Name,
			CreatedAt:   timestamppb.New(order.ShopCard.CreatedAt),
			Description: order.ShopCard.Description,
			Category:    order.ShopCard.Category,
			Price:       *order.ShopCard.Price,
			Visible:     order.ShopCard.Visible,
			CoverUrl:    *order.ShopCard.CoverURL,
		},
		Price: order.Price,
	}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := s.db.GetOrder(
		req.OrderUuid,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}
	return &pb.OrderResponse{
		Uuid:         order.UUID,
		CustomerUuid: order.CustomerUUID,
		Status:       string(order.Status),
		ShopCard: &pb.ShopCard{
			Uuid:        order.ShopCard.UUID,
			Name:        order.ShopCard.Name,
			CreatedAt:   timestamppb.New(order.ShopCard.CreatedAt),
			Description: order.ShopCard.Description,
			Category:    order.ShopCard.Category,
			Price:       *order.ShopCard.Price,
			Visible:     order.ShopCard.Visible,
			CoverUrl:    *order.ShopCard.CoverURL,
		},
		Price: order.Price,
	}, nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

func isValidStatus(status string) bool {
	switch db.OrderStatus(status) {
	case db.StatusNew,
		db.StatusUnpaid,
		db.StatusPaid,
		db.StatusCancelled,
		db.StatusRefunded:
		return true
	default:
		return false
	}
}
