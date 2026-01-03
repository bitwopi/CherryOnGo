package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	pb "shopcards/server/api/grpc/shop_card"
	"shopcards/server/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedShopCardServiceServer
	db db.DBManager
}

func NewServer(dsn string) *Server {
	dbManager, err := db.NewManager(dsn)
	if err != nil {
		log.Fatalf("failed to create db manager: %v", err)
	}
	return &Server{db: dbManager}
}

func (s *Server) Start(bindURL string) error {
	log.Println("Starting gRPC server...")
	s.db.Migrate()
	server := grpc.NewServer()
	pb.RegisterShopCardServiceServer(server, s)
	lis, err := net.Listen("tcp", bindURL)
	if err != nil {
		panic(err)
	}
	log.Println("gRPC server listening on: ", bindURL)
	return server.Serve(lis)
}

func (s *Server) CreateShopCard(ctx context.Context, req *pb.ShopCardRequest) (*pb.ShopCardResponse, error) {
	card, err := s.db.CreateShopCard(req.Name, req.Description, req.Category, &req.Price, req.Visible, &req.CoverUrl)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.AlreadyExists, "card is already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create card")
	}
	return &pb.ShopCardResponse{
		Uuid:        card.UUID,
		Name:        card.Name,
		Description: card.Description,
		Category:    card.Category,
		Price:       *card.Price,
		Visible:     card.Visible,
		CoverUrl:    *card.CoverURL,
		CreatedAt:   timestamppb.New(card.CreatedAt),
	}, nil
}

func (s *Server) UpdateShopCard(ctx context.Context, req *pb.UpdateShopCardRequest) (*pb.ShopCardResponse, error) {
	card, err := s.db.UpdateShopCard(req.Uuid, req.Name, req.Description, req.Category, &req.Price, req.Visible, &req.CoverUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update card")
	}
	return &pb.ShopCardResponse{
		Uuid:        card.UUID,
		Name:        card.Name,
		Description: card.Description,
		Category:    card.Category,
		Price:       *card.Price,
		Visible:     card.Visible,
		CoverUrl:    *card.CoverURL,
		CreatedAt:   timestamppb.New(card.CreatedAt),
	}, nil
}

func (s *Server) GetShopCard(ctx context.Context, req *pb.ShopCardUUIDRequest) (*pb.ShopCardResponse, error) {
	card, err := s.db.GetShopCard(req.ShopCardUuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "card not found")
	}
	return &pb.ShopCardResponse{
		Uuid:        card.UUID,
		Name:        card.Name,
		Description: card.Description,
		Category:    card.Category,
		Price:       *card.Price,
		Visible:     card.Visible,
		CoverUrl:    *card.CoverURL,
		CreatedAt:   timestamppb.New(card.CreatedAt),
	}, nil
}

func (s *Server) DeleteShopCard(ctx context.Context, req *pb.ShopCardUUIDRequest) (*emptypb.Empty, error) {
	err := s.db.DeleteShopCard(req.ShopCardUuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "card not found")
	}
	return &emptypb.Empty{}, nil
}
