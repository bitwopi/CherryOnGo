package tests

import (
	"testing"
	users "users/server/api/grpc/user"
	"users/tests/suites"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	passDeafaultLen = 10
)

func TestRegisterLoginRefreshPositive(t *testing.T) {
	ctx, st := suites.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, false, false, passDeafaultLen)
	authRequest := &users.AuthRequest{
		Login:    email,
		Password: pass,
	}

	regResponse, err := st.UserClient.SignUpUser(ctx, authRequest)
	require.NoError(t, err)
	assert.Equal(t, "user created", regResponse.Status)

	authResponse, err := st.UserClient.AuthUser(ctx, authRequest)
	require.NoError(t, err)
	assert.NotEmpty(t, authResponse.AccessToken)
	assert.NotEmpty(t, authResponse.RefreshToken)

	parsedToken, err := st.JWTManager.ParseJWT(authResponse.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, regResponse.UserUuid, parsedToken.Subject)

	refreshResponse, err := st.UserClient.RefreshJWT(
		ctx,
		&users.RefreshRequest{
			RefreshToken: authResponse.RefreshToken,
		},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshResponse.AccessToken)
	assert.NotEmpty(t, refreshResponse.RefreshToken)

	refreshResponse, err = st.UserClient.RefreshJWT(
		ctx,
		&users.RefreshRequest{
			RefreshToken: authResponse.RefreshToken,
		},
	)
	require.Error(t, err)
}
