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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RemnaGRPCServer struct {
	pb.UnimplementedRemnaServiceServer

	remnaClient client.Client
	grpcServer  *grpc.Server
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
		panic(err)
	}
	r.grpcServer = server

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
		return nil, status.Error(codes.NotFound, err.Error())
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
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := convertUserResponse(resp.Response)
	return result, nil
}

func (r *RemnaGRPCServer) UpdateUserExpiryTime(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	log.Println("UpdateUserExpiryTime called", "req", req)
	plan := client.Plans[req.Plan]
	if plan == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid plan")
	}
	resp, err := r.remnaClient.UpdateUserExpiryTime(ctx, plan, req.Username, req.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := convertUserResponse(resp.Response)
	log.Println(result)
	return result, nil
}

func (r *RemnaGRPCServer) GetUsersByTgID(ctx context.Context, req *pb.GetUserByTgIDRequest) (*pb.MultipleUsersResponse, error) {
	resp, err := r.remnaClient.GetUsersByTgID(ctx, req.Tgid)
	log.Println("GetUsersByTgID response: ", resp)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertMultipleUsersResponse(*resp), nil
}

func (r *RemnaGRPCServer) GetUsersByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.MultipleUsersResponse, error) {
	resp, err := r.remnaClient.GetUsersByEmail(ctx, req.Email)
	log.Println("GetUsersByEmail response: ", resp)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertMultipleUsersResponse(*resp), nil
}

func (r *RemnaGRPCServer) GetAllUsers(ctx context.Context, req *pb.EmptyRequest) (*pb.MultipleUsersResponse, error) {
	resp, err := r.remnaClient.GetAllUsers(ctx)
	log.Println("GetAllUsers response: ", resp)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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

func (r *RemnaGRPCServer) DisableUser(ctx context.Context, req *pb.GetUserUUIDRequest) (*pb.UserResponse, error) {
	resp, err := r.remnaClient.DisableUser(ctx, req.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to disable user")
	}

	result := convertUserResponse(resp.Response)
	log.Println(result)

	return convertUserResponse(resp.Response), nil
}

func (r *RemnaGRPCServer) EnableUser(ctx context.Context, req *pb.GetUserUUIDRequest) (*pb.UserResponse, error) {
	resp, err := r.remnaClient.EnableUser(ctx, req.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to enable user")
	}

	return convertUserResponse(resp.Response), nil
}

func (r *RemnaGRPCServer) GetUserHwidDevices(ctx context.Context, req *pb.GetUserUUIDRequest) (*pb.MultipleHwidResponse, error) {
	resp, err := r.remnaClient.GetUserHwidDevices(ctx, req.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	var devices []*pb.HwidResponse
	for _, d := range resp.Response.Devices {
		devices = append(devices, &pb.HwidResponse{
			Hwid:        d.Hwid,
			UserUuid:    d.UserUuid.String(),
			Platform:    d.Platform.Value,
			OsVersion:   d.OsVersion.Value,
			DeviceModel: d.DeviceModel.Value,
			UserAgent:   d.UserAgent.Value,
			CreatedAt:   timestamppb.New(d.CreatedAt),
			UpdatedAt:   timestamppb.New(d.UpdatedAt),
		})
	}
	return &pb.MultipleHwidResponse{
		Devices: devices,
	}, nil
}

func (r *RemnaGRPCServer) GetSRHHistory(ctx context.Context, req *pb.EmptyRequest) (*pb.SRHHistoryResponse, error) {
	resp, err := r.remnaClient.GetSRHHistory(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	var records []*pb.SRHRecordResponse
	for _, record := range resp.Records {
		records = append(records, &pb.SRHRecordResponse{
			Id:          int64(record.ID),
			UserUuid:    record.UserUuid.String(),
			RequestIp:   record.RequestIp.Value,
			UserAgent:   record.UserAgent.Value,
			RequestedAt: timestamppb.New(record.RequestAt),
		})
	}
	return &pb.SRHHistoryResponse{
		Records: records,
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

func (r *RemnaGRPCServer) Stop() {
	r.grpcServer.GracefulStop()
}
