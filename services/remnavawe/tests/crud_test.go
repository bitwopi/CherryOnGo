package tests

import (
	pb "remnawave/server/api/grpc/remna"
	"remnawave/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestCRUDPositive(t *testing.T) {
	ctx, st := suite.New(t)
	username := "test_user_client"
	plan := pb.Plan{
		DeviceLimit: 3,
		DayLimit:    30,
		Squad:       "basic",
	}
	req := pb.CreateUserRequest{
		Username: username,
		Email:    gofakeit.Email(),
		Plan:     &plan,
	}

	resp, err := st.RemnaClient.CreateUser(ctx, &req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.EqualValues(t, plan.DeviceLimit, resp.DeviceLimit)
	getResp, err := st.RemnaClient.GetUser(ctx, &pb.GetUserByUsernameRequest{Username: username})
	assert.NoError(t, err)
	assert.NotEmpty(t, getResp)
	assert.EqualValues(t, resp.Uuid, getResp.Uuid)
	updReq := pb.UpdateUserRequest{
		Uuid:     resp.Uuid,
		Username: username,
		Plan:     &plan,
	}
	updResp, err := st.RemnaClient.UpdateUserExpiryTime(ctx, &updReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, updResp)
	assert.NotEqual(t, getResp.ExpiryTime.AsTime(), updResp.ExpiryTime.AsTime())
}
