package grpc

import (
	"context"
	"log"
	pb "orders/server/api/grpc/remna"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Address        string
	RequestTimeout time.Duration
	MaxRetries     uint
}

type RemnaGRPCClient struct {
	conn   *grpc.ClientConn
	client pb.RemnaServiceClient
	timout time.Duration
}

func NewRemnaGRPCClient(cfg Config) (*RemnaGRPCClient, error) {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(cfg.MaxRetries),
		grpc_retry.WithPerRetryTimeout(500 * time.Millisecond),
		grpc_retry.WithCodes(
			codes.DeadlineExceeded,
			codes.Unavailable,
		),
	}

	conn, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewRemnaServiceClient(conn)
	return &RemnaGRPCClient{
		conn:   conn,
		client: client,
		timout: cfg.RequestTimeout,
	}, nil
}

func (c *RemnaGRPCClient) Close() error {
	return c.conn.Close()
}

func (c *RemnaGRPCClient) GetUser(username string) *pb.UserResponse {
	ctx, cancel := context.WithTimeout(context.Background(), c.timout)
	defer cancel()
	req := &pb.GetUserRequest{
		Username: username,
	}
	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		log.Println(err)
		return nil
	}
	return resp
}

func (c *RemnaGRPCClient) CreateUser(username string) *pb.UserResponse {
	ctx, cancel := context.WithTimeout(context.Background(), c.timout)
	defer cancel()
	req := &pb.GetUserRequest{
		Username: username,
	}
	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil
	}
	return resp
}
