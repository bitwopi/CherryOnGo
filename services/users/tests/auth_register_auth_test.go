package tests

import (
	"testing"
	users "users/server/api/grpc/users"
	suites "users/tests/suite"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	passDeafaultLen = 10
)

func TestSignUpSignInRefreshPositive(t *testing.T) {
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

func TestSignUpNegative(t *testing.T) {
	ctx, st := suites.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, false, false, passDeafaultLen)
	authRequestNoPass := &users.AuthRequest{
		Login:    email,
		Password: "",
	}
	authRequestNoEmail := &users.AuthRequest{
		Login:    "",
		Password: pass,
	}

	regResponse, err := st.UserClient.SignUpUser(ctx, authRequestNoPass)
	require.ErrorContains(t, err, "empty login or password")
	assert.Empty(t, regResponse)
	regResponse, err = st.UserClient.SignUpUser(ctx, authRequestNoEmail)
	require.ErrorContains(t, err, "empty login or password")
	assert.Empty(t, regResponse)
}

func TestSignInNegative(t *testing.T) {
	ctx, st := suites.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, false, false, passDeafaultLen)
	authRequest := &users.AuthRequest{
		Login:    email,
		Password: pass,
	}
	authRequestNoPass := &users.AuthRequest{
		Login:    email,
		Password: "",
	}
	authRequestNoEmail := &users.AuthRequest{
		Login:    "",
		Password: pass,
	}
	authRequestWrongEmail := &users.AuthRequest{
		Login:    "aboba@gmail.com",
		Password: pass,
	}
	authRequestWrongPass := &users.AuthRequest{
		Login:    email,
		Password: "PisyaPopa221",
	}

	_, err := st.UserClient.SignUpUser(ctx, authRequest)
	require.NoError(t, err)

	authResponse, err := st.UserClient.AuthUser(ctx, authRequestNoPass)
	require.ErrorContains(t, err, "empty login or password")
	assert.Empty(t, authResponse)

	authResponse, err = st.UserClient.AuthUser(ctx, authRequestNoEmail)
	require.ErrorContains(t, err, "empty login or password")
	assert.Empty(t, authResponse)

	authResponse, err = st.UserClient.AuthUser(ctx, authRequestWrongEmail)
	require.ErrorContains(t, err, "failed to get user")
	assert.Empty(t, authResponse)

	authResponse, err = st.UserClient.AuthUser(ctx, authRequestWrongPass)
	require.ErrorContains(t, err, "invalid password")
	assert.Empty(t, authResponse)

}
