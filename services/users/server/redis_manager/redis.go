package redismanager

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	client *redis.Client
}

type TokenSession struct {
	UserID       string        `redis:"user_id"`
	AccessToken  string        `redis:"access_token"`
	RefreshToken string        `redis:"refresh_token"`
	CreatedAt    time.Time     `redis:"created_at"`
	TTL          time.Duration `redis:"ttl"`
	ExpiresAt    time.Time     `redis:"expires_at"`
}

func NewRedisManager(addr string, pass string, dbnum int) *RedisManager {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       dbnum,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return &RedisManager{client: client}
}

func (m *RedisManager) CreateSession(aToken string, rToken string, userID string, ttl time.Duration) (*TokenSession, error) {
	createdAt := time.Now()
	session := TokenSession{
		UserID:       userID,
		AccessToken:  aToken,
		RefreshToken: rToken,
		CreatedAt:    createdAt,
		TTL:          ttl,
		ExpiresAt:    createdAt.Add(ttl),
	}

	key := fmt.Sprintf("session:%s", rToken)
	ctx := context.Background()
	if err := m.client.HSet(ctx, key, session).Err(); err != nil {
		return nil, err
	}
	m.client.Expire(ctx, key, ttl)
	return &session, nil
}

func (m *RedisManager) GetSession(rToken string) (*TokenSession, error) {
	key := fmt.Sprintf("session:%s", rToken)
	ctx := context.Background()
	var session TokenSession
	if err := m.client.HGetAll(ctx, key).Scan(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (m *RedisManager) DeleteSession(rToken string) error {
	key := fmt.Sprintf("session:%s", rToken)
	ctx := context.Background()
	if err := m.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
