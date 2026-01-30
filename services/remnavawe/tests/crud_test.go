package tests

import (
	"remnawave/client"
	pb "remnawave/server/api/grpc/remna"
	"remnawave/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCRUDPositive(t *testing.T) {
	ctx, st := suite.New(t)
	username := "test_user_client4"
	req := pb.CreateUserRequest{
		Username: username,
		Email:    gofakeit.Email(),
		Plan:     "3:30",
	}

	resp, err := st.RemnaClient.CreateUser(ctx, &req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.EqualValues(t, client.Plans["3:30"].DeviceLimit, resp.DeviceLimit)
	getResp, err := st.RemnaClient.GetUser(ctx, &pb.GetUserRequest{Username: username})
	assert.NoError(t, err)
	assert.NotEmpty(t, getResp)
	assert.EqualValues(t, resp.Uuid, getResp.Uuid)
	updReq := pb.UpdateUserRequest{
		Uuid:     resp.Uuid,
		Username: username,
		Plan:     "3:30",
	}
	updResp, err := st.RemnaClient.UpdateUserExpiryTime(ctx, &updReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, updResp)
	assert.NotEqual(t, getResp.ExpiryTime, updResp.ExpiryTime)
}
