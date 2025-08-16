package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error)
	GetUserID(ctx context.Context, session string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
