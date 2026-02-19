package remnaclient

import (
	"context"
	pb "gateway/server/api/gateway/grpc/gen/remna"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type RemnaGRPCClient struct {
	conn    *grpc.ClientConn
	client  pb.RemnaServiceClient
	timeout time.Duration
}

func NewRemnaGRPCClient(addr string, timeout time.Duration, maxRetries uint) (*RemnaGRPCClient, error) {
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
	client := pb.NewRemnaServiceClient(conn)
	return &RemnaGRPCClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *RemnaGRPCClient) Close() error {
	return c.conn.Close()
}

func (c *RemnaGRPCClient) GetUsersByEmail(email string) (*pb.MultipleUsersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.GetUserByEmailRequest{
		Email: email,
	}
	resp, err := c.client.GetUsersByEmail(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
