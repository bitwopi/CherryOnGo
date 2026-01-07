package tests

import (
	pb "orders/server/api/grpc/gen/order"
	"orders/server/db"
	"orders/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCRUDPositive(t *testing.T) {
	ctx, st := suite.New(t)
	cRequest := pb.CreateOrderRequest{
		CustomerUuid: uuid.NewString(),
		Status:       string(db.StatusNew),
		ShopCard: &pb.ShopCard{
			Uuid:        uuid.NewString(),
			Name:        gofakeit.BookTitle(),
			Description: gofakeit.Product().Description,
			Category:    gofakeit.ProductCategory(),
			CreatedAt:   timestamppb.Now(),
			Visible:     true,
		},
		Price: gofakeit.Float32(),
	}

	resp, err := st.OrderClient.CreateOrder(ctx, &cRequest)
	require.NoError(t, err)
	assert.NotEmpty(t, resp)

	uRequest := pb.OrderStatusRequest{
		OrderUuid: resp.Uuid,
		Status:    string(db.StatusUnpaid),
	}
	uResp, err := st.OrderClient.UpdateOrderStatus(ctx, &uRequest)
	require.NoError(t, err)
	assert.NotEmpty(t, uResp)

	gRequest := pb.GetOrderRequest{OrderUuid: resp.Uuid}
	resp, err = st.OrderClient.GetOrder(ctx, &gRequest)
	require.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, uResp.Status, resp.Status)
}
