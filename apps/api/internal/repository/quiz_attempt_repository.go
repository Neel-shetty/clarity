package repository

import (
	"context"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/gofrs/uuid/v5"
)

type QuizAttemptRepository interface {
	CreateAttempt(ctx context.Context, attempt *model.QuizAttempt) error
	UpdateAttempt(ctx context.Context, attempt *model.QuizAttempt) error
	FindAttempByUser(ctx context.Context, userID uuid.UUID) ([]*model.QuizAttempt, error)
	FindAttempByQuiz(ctx context.Context, quizID uuid.UUID) ([]*model.QuizAttempt, error)
}
