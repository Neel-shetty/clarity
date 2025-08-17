package repository

import (
	"context"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/gofrs/uuid/v5"
)

type QuizRepository interface {
	CreateQuiz(ctx context.Context, quiz *model.Quiz) error
	UpdateQuiz(ctx context.Context, quiz *model.Quiz) error
	DeleteQuiz(ctx context.Context, quizID uuid.UUID) error
	FindbyID(ctx context.Context, quizID uuid.UUID) (*model.Quiz, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Quiz, error)

	AddQuestion(ctx context.Context, question *model.Question) error
	UpdateQuestion(ctx context.Context, question *model.Question) error
	DeleteQuestion(ctx context.Context, questionID uuid.UUID) error
}
