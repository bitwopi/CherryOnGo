package shopcardclient

import (
	"context"
	"fmt"
	"time"

	pb "gateway/server/api/gateway/grpc/gen/shop_card"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type ShopCardGRPCClient struct {
	conn    *grpc.ClientConn
	client  pb.ShopCardServiceClient
	timeout time.Duration
}

func NewShopCardGRPCClient(addr string, timeout time.Duration, maxRetries uint) (*ShopCardGRPCClient, error) {
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
	client := pb.NewShopCardServiceClient(conn)
	return &ShopCardGRPCClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *ShopCardGRPCClient) Close() error {
	return c.conn.Close()
}

func (c *ShopCardGRPCClient) CreateShopCard(
	name string,
	description string,
	category string,
	price float32,
	visible bool,
	coverUrl string,
	metadata map[string]interface{}) (*pb.ShopCardResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	metaStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata %v", err)
	}
	req := &pb.ShopCardRequest{
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		Visible:     visible,
		CoverUrl:    coverUrl,
		Metadata:    metaStruct,
	}
	resp, err := c.client.CreateShopCard(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ShopCardGRPCClient) UpdateShopCard(
	uuid string,
	name string,
	description string,
	category string,
	price float32,
	visible bool,
	coverUrl string,
	metadata map[string]interface{}) (*pb.ShopCardResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	metaStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata %v", err)
	}
	req := &pb.UpdateShopCardRequest{
		Uuid:        uuid,
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		Visible:     visible,
		CoverUrl:    coverUrl,
		Metadata:    metaStruct,
	}
	resp, err := c.client.UpdateShopCard(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ShopCardGRPCClient) GetShopCard(uuid string) (*pb.ShopCardResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.ShopCardUUIDRequest{
		ShopCardUuid: uuid,
	}
	resp, err := c.client.GetShopCard(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ShopCardGRPCClient) GetAllShopCards() (*pb.MultipleResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.client.GetAllShopCards(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ShopCardGRPCClient) DeleteShopCard(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req := &pb.ShopCardUUIDRequest{
		ShopCardUuid: uuid,
	}
	_, err := c.client.DeleteShopCard(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
