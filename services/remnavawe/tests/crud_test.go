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
	username := "test_user_client3"
	plan := pb.Plan{
		DeviceLimit:       3,
		DayLimit:          30,
		TrafficLimitBytes: 10 * 1024 * 1024 * 1024,
		Squad:             "f0bb8401-22ee-4b67-b256-d24cd64ee102",
	}
	tgID := "123456789"
	req := pb.CreateUserRequest{
		Username: username,
		Email:    gofakeit.Email(),
		Plan:     &plan,
		Tgid:     tgID,
	}

	resp, err := st.RemnaClient.CreateUser(ctx, &req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.EqualValues(t, plan.DeviceLimit, resp.DeviceLimit)
	assert.EqualValues(t, plan.TrafficLimitBytes, resp.TrafficLimitBytes)
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
	updResp, err = st.RemnaClient.AddUserTraffic(ctx, &updReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, updResp)
	assert.NotEqual(t, getResp.TrafficLimitBytes, updResp.TrafficLimitBytes)
}
