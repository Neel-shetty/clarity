package service

import (
	"context"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/Neel-shetty/clarity/internal/repository"
	"github.com/gofrs/uuid/v5"
)

type QuizService interface {
	CreateQuiz(ctx context.Context, userID uuid.UUID, title string, noteIDs []*model.Note) (*model.Quiz, error)
	GetQuizByID(ctx context.Context, id uuid.UUID) (*model.Quiz, error)
	ListUserQuizzes(ctx context.Context, userID uuid.UUID) ([]*model.Quiz, error)
	DeleteQuiz(ctx context.Context, id uuid.UUID) error
}

type quizService struct {
	repo repository.QuizRepository
}

func NewQuizService(repo repository.QuizRepository) QuizService {
	return &quizService{repo}
}

func (s *quizService) CreateQuiz(ctx context.Context, userID uuid.UUID, title string, notes []*model.Note) (*model.Quiz, error) {
	uuid, err:= uuid.NewV7()
	if err != nil {
		return nil, err
	}
	quiz := &model.Quiz{
		ID: uuid,
		UserID: userID,
		Notes:  notes,
		Title:  title,
	}
	if err := s.repo.CreateQuiz(ctx, quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *quizService) GetQuizByID(ctx context.Context, id uuid.UUID) (*model.Quiz, error) {
	quiz, err := s.repo.FindbyID(ctx, id)
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *quizService) ListUserQuizzes(ctx context.Context, userID uuid.UUID) ([]*model.Quiz, error) {
	quiz, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *quizService) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteQuiz(ctx, id); err != nil {
		return err
	}
	return nil
}
