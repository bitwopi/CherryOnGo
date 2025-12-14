package grpc

import (
	"context"
	"log"
	"net"
	"remnawave/client"
	"remnawave/config"
	pb "remnawave/server/api/grpc/remna"
	"strconv"

	"github.com/Jolymmiles/remnawave-api-go/v2/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RemnaGRPCServer struct {
	pb.UnimplementedRemnaServiceServer

	remnaClient client.Client
}

func NewRemnaGRPCServer(cfg *config.Config) *RemnaGRPCServer {
	client := client.NewClient(cfg.RemnaAPIKey, cfg.RemnaURL)
	return &RemnaGRPCServer{
		remnaClient: *client,
	}
}

func (r *RemnaGRPCServer) Start(bindURL string) error {
	log.Println("Starting gRPC server...")
	server := grpc.NewServer()
	pb.RegisterRemnaServiceServer(server, r)
	lis, err := net.Listen("tcp", bindURL)
	if err != nil {
		return err
	}

	log.Println("gRPC server listening on: ", bindURL)
	return server.Serve(lis)
}

func (r *RemnaGRPCServer) PingRemna(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	err := r.remnaClient.Ping()
	if err != nil {
		return &pb.PingResponse{
			Status: "down",
		}, err
	}
	return &pb.PingResponse{
		Status: "up",
	}, nil
}
func (r *RemnaGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	log.Println("GetUser called", "req", req)
	defer log.Println("GetUser finished")
	resp, err := r.remnaClient.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	result := convertUserResponse(*resp)
	return result, nil

}
func (r *RemnaGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Println("CreateUser called", "req", req)
	defer log.Println("CreateUser finished")
	resp, err := r.remnaClient.CreateUser(client.Plans[req.Plan], req.Username, req.Tgid, req.Email)
	if err != nil {
		return nil, err
	}
	result := convertUserResponse(*resp)
	return result, nil
}

func convertUserResponse(resp api.UserResponse) *pb.UserResponse {
	ur := resp.Response
	var intSquad []string
	for _, squad := range ur.ActiveInternalSquads {
		intSquad = append(intSquad, squad.GetUUID().String())
	}
	result := pb.UserResponse{
		Uuid:           ur.UUID.String(),
		Username:       ur.Username,
		Tgid:           strconv.Itoa(ur.TelegramId.Value),
		Email:          ur.Email.Value,
		InternalSquads: intSquad,
		ExpiryTime:     timestamppb.New(ur.GetExpireAt()),
		SubUrl:         ur.SubscriptionUrl,
	}
	return &result
}
