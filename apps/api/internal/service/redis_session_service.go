package service

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
)

type RedisSessionService struct {
	client *redis.Client
	prefix string
}

func NewRedisSessionService(client *redis.Client) *RedisSessionService {
	return &RedisSessionService{
		client: client,
		prefix: "session:",
	}
}

func (s *RedisSessionService) CreateSession(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	sessionID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	key := s.prefix + sessionID.String()
	err = s.client.Set(ctx, key, userID.String(), ttl).Err()
	if err != nil {
		return "", err
	}
	return sessionID.String(), nil
}

func (s *RedisSessionService) GetUserID(ctx context.Context, sessionID string) (uuid.UUID, error) {
	val, err := s.client.Get(ctx, s.prefix+sessionID).Result()
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.FromString(val)
}

func (s *RedisSessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return s.client.Del(ctx, s.prefix+sessionID).Err()
}
