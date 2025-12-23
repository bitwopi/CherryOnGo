package jwtmanager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	signingKey string
}

func NewJWTManager(signingKey string) (*JWTManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JWTManager{signingKey: signingKey}, nil
}

func (m *JWTManager) NewJWT(userID string, ttl time.Duration) (string, error) {
	if len(userID) != 36 || ttl == 0 {
		return "", errors.New("invalid token data")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	})
	return token.SignedString(m.signingKey)
}

func (m *JWTManager) NewRefreshToken() string {
	return uuid.NewString()
}

func (m *JWTManager) ParseJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected singning method")
		}
		return m.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
