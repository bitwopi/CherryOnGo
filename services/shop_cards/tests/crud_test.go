package tests

import (
	pb "shopcards/server/api/grpc/shop_card"
	"shopcards/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCRUDPositive(t *testing.T) {
	ctx, st := suite.New(t)
	name := gofakeit.BookTitle()
	desc := gofakeit.Product().Description
	category := "other"
	price := float32(gofakeit.Price(0, 40000))
	coverUrl := gofakeit.URL()
	req := pb.ShopCardRequest{
		Name:        name,
		Description: desc,
		Category:    category,
		Price:       price,
		CoverUrl:    coverUrl,
	}
	resp, err := st.ShopCardClient.CreateShopCard(ctx, &req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp)
	getResp, err := st.ShopCardClient.GetShopCard(ctx, &pb.ShopCardUUIDRequest{ShopCardUuid: resp.Uuid})
	require.NoError(t, err)
	assert.NotEmpty(t, getResp)
	newName := gofakeit.BookTitle()
	updReq := pb.UpdateShopCardRequest{
		Uuid:        resp.Uuid,
		Name:        newName,
		Description: desc,
		Category:    category,
		Price:       price,
		CoverUrl:    coverUrl,
		Visible:     true,
	}
	updResp, err := st.ShopCardClient.UpdateShopCard(ctx, &updReq)
	require.NoError(t, err)
	assert.NotEmpty(t, updResp)
	assert.Equal(t, newName, updResp.Name)
	_, err = st.ShopCardClient.DeleteShopCard(ctx, &pb.ShopCardUUIDRequest{ShopCardUuid: resp.Uuid})
	require.NoError(t, err)
}
