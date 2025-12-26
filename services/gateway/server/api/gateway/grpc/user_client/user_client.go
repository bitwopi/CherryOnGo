package userclient

import (
	"context"
	pb "gateway/server/api/gateway/grpc/gen/users"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGRPCClient struct {
	conn    *grpc.ClientConn
	client  pb.UserServiceClient
	timeout time.Duration
}

func NewUserGRPCClient(addr string, timeout time.Duration, maxRetries uint) (*UserGRPCClient, error) {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(maxRetries),
		grpc_retry.WithPerRetryTimeout(500 * time.Millisecond),
		grpc_retry.WithCodes(
			codes.DeadlineExceeded,
			codes.Unavailable,
		),
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewUserServiceClient(conn)
	return &UserGRPCClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *UserGRPCClient) Close() error {
	return c.conn.Close()
}

func (c *UserGRPCClient) SignUpUser(login string, password string) (*pb.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.AuthRequest{
		Login:    login,
		Password: password,
	}
	resp, err := c.client.SignUpUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *UserGRPCClient) AuthUser(login string, password string) (*pb.JWTResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.AuthRequest{
		Login:    login,
		Password: password,
	}
	resp, err := c.client.AuthUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *UserGRPCClient) RefreshJWT(refreshToken string) (*pb.JWTResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.RefreshRequest{
		RefreshToken: refreshToken,
	}
	resp, err := c.client.RefreshJWT(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
