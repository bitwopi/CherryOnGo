package tests

import (
	pb "remnawave/server/api/grpc/remna"
	"remnawave/tests/suite"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	ctx, st := suite.New(t)
	req := pb.EmptyRequest{}
	resp, err := st.RemnaClient.GetAllUsers(ctx, &req)
	assert.NoError(t, err, "GetAllUsers failed")
	assert.NotNil(t, resp, "GetAllUsers response is nil")
	assert.Greater(t, len(resp.Users), 0, "No users returned")
}

func TestGetUsersWithParameters(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	tg := strconv.FormatInt(gofakeit.Int64(), 10)
	req := pb.CreateUserRequest{
		Username: gofakeit.Username(),
		Email:    email,
		Plan:     "3:30",
		Tgid:     tg,
	}

	resp, err := st.RemnaClient.CreateUser(ctx, &req)
	assert.NoError(t, err, "CreateUser failed")
	assert.NotNil(t, resp, "CreateUser response is nil")

	eReq := pb.GetUserByEmailRequest{Email: email}
	eResp, err := st.RemnaClient.GetUsersByEmail(ctx, &eReq)
	assert.NoError(t, err, "GetUsersByEmail failed")
	assert.NotNil(t, eResp, "GetUsersByEmail response is nil")
	assert.Greater(t, len(eResp.Users), 0, "No users returned by email")
	assert.EqualValues(t, email, eResp.Users[0].Email, "Email does not match")

	tgReq := pb.GetUserByTgIDRequest{Tgid: tg}
	tgResp, err := st.RemnaClient.GetUsersByTgID(ctx, &tgReq)
	assert.NoError(t, err, "GetUsersByTgID failed")
	assert.NotNil(t, tgResp, "GetUsersByTgID response is nil")
	assert.Greater(t, len(tgResp.Users), 0, "No users returned by tg id")
	assert.EqualValues(t, tg, tgResp.Users[0].Tgid, "TgID does not match")
}
