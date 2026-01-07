package orderclient

import (
	"context"
	pb "gateway/server/api/gateway/grpc/gen/order"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderGRPCClient struct {
	conn    *grpc.ClientConn
	client  pb.OrderServiceClient
	timeout time.Duration
}

func NewOrderGRPCClient(addr string, timeout time.Duration, maxRetries uint) (*OrderGRPCClient, error) {
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
	client := pb.NewOrderServiceClient(conn)
	return &OrderGRPCClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *OrderGRPCClient) Close() error {
	return c.conn.Close()
}

func (c *OrderGRPCClient) CreateOrder(customerUUID string, status string, shopCard *pb.ShopCard, price float32) (*pb.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.CreateOrderRequest{
		CustomerUuid: customerUUID,
		Status:       status,
		ShopCard:     shopCard,
		Price:        price,
	}
	resp, err := c.client.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *OrderGRPCClient) GetOrder(orderUUID string) (*pb.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.GetOrderRequest{
		OrderUuid: orderUUID,
	}
	resp, err := c.client.GetOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *OrderGRPCClient) UpdateOrderStatus(orderUUID string, status string) (*pb.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.OrderStatusRequest{
		OrderUuid: orderUUID,
		Status:    status,
	}
	resp, err := c.client.UpdateOrderStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
