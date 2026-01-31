package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"remnawave/client"
	"remnawave/config"
	pb "remnawave/server/api/grpc/remna"
	"strconv"

	"github.com/Jolymmiles/remnawave-api-go/v2/api"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RemnaGRPCServer struct {
	pb.UnimplementedRemnaServiceServer

	remnaClient client.Client
}

func NewRemnaGRPCServer(cfg *config.Config) *RemnaGRPCServer {
	log.Println(cfg.Remna)
	client := client.NewClient(cfg.Remna.RemnaAPIKey, cfg.Remna.RemnaURL)
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

func (r *RemnaGRPCServer) PingRemna(ctx context.Context, req *pb.EmptyRequest) (*pb.PingResponse, error) {
	err := r.remnaClient.Ping(ctx)
	if err != nil {
		return &pb.PingResponse{
			Status: "down",
		}, err
	}
	return &pb.PingResponse{
		Status: "up",
	}, nil
}
func (r *RemnaGRPCServer) GetUser(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.UserResponse, error) {
	log.Println("GetUser called", "req", req)
	defer log.Println("GetUser finished")
	resp, err := r.remnaClient.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	log.Println(resp)
	result := convertUserResponse(resp.Response)
	return result, nil

}
func (r *RemnaGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Println("CreateUser called", "req", req)
	defer log.Println("CreateUser finished")
	resp, err := r.remnaClient.CreateUser(ctx, client.Plans[req.Plan], req.Username, req.Tgid, req.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := convertUserResponse(resp.Response)
	return result, nil
}

func (r *RemnaGRPCServer) UpdateUserExpiryTime(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	userUUID, err := uuid.Parse(req.Uuid)
	plan := client.Plans[req.Plan]
	if plan == nil {
		return nil, errors.New("invalid plan")
	}
	resp, err := r.remnaClient.UpdateUserExpiryTime(ctx, plan, &req.Username, &userUUID)
	if err != nil {
		return nil, err
	}
	result := convertUserResponse(resp.Response)
	return result, nil
}

func (r *RemnaGRPCServer) GetUsersByTgID(ctx context.Context, req *pb.GetUserByTgIDRequest) (*pb.MultipleUsersResponse, error) {
	resp, err := r.remnaClient.GetUsersByTgID(ctx, req.Tgid)
	log.Println("GetUsersByTgID response: ", resp)
	if err != nil {
		return nil, err
	}
	return convertMultipleUsersResponse(*resp), nil
}

func (r *RemnaGRPCServer) GetUsersByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.MultipleUsersResponse, error) {
	resp, err := r.remnaClient.GetUsersByEmail(ctx, req.Email)
	log.Println("GetUsersByEmail response: ", resp)
	if err != nil {
		return nil, err
	}
	return convertMultipleUsersResponse(*resp), nil
}

func (r *RemnaGRPCServer) GetAllUsers(ctx context.Context, req *pb.EmptyRequest) (*pb.MultipleUsersResponse, error) {
	resp, err := r.remnaClient.GetAllUsers(ctx)
	log.Println("GetAllUsers response: ", resp)
	if err != nil {
		return nil, err
	}

	var users []*pb.UserResponse
	for _, u := range resp.Response.Users {
		var intSquad []string
		for _, squad := range u.ActiveInternalSquads {
			intSquad = append(intSquad, squad.GetUUID().String())
		}
		result := pb.UserResponse{
			Uuid:           u.UUID.String(),
			Username:       u.Username,
			Tgid:           strconv.Itoa(u.TelegramId.Value),
			Email:          u.Email.Value,
			InternalSquads: intSquad,
			ExpiryTime:     timestamppb.New(u.GetExpireAt()),
			SubUrl:         u.SubscriptionUrl,
			DeviceLimit:    int64(u.HwidDeviceLimit.Value),
		}
		users = append(users, &result)
	}
	log.Println(users)
	return &pb.MultipleUsersResponse{
		Users: users,
	}, nil
}
func convertUserResponse(resp api.User) *pb.UserResponse {
	var intSquad []string
	for _, squad := range resp.ActiveInternalSquads {
		intSquad = append(intSquad, squad.GetUUID().String())
	}
	result := pb.UserResponse{
		Uuid:           resp.UUID.String(),
		Username:       resp.Username,
		Tgid:           strconv.Itoa(resp.TelegramId.Value),
		Email:          resp.Email.Value,
		InternalSquads: intSquad,
		ExpiryTime:     timestamppb.New(resp.GetExpireAt()),
		SubUrl:         resp.SubscriptionUrl,
		DeviceLimit:    int64(resp.HwidDeviceLimit.Value),
	}
	return &result
}

func convertMultipleUsersResponse(resp api.UsersResponse) *pb.MultipleUsersResponse {
	ur := resp.Response
	var users []*pb.UserResponse
	for _, u := range ur {
		users = append(users, convertUserResponse(u))
	}
	return &pb.MultipleUsersResponse{
		Users: users,
	}
}
