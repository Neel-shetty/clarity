package repository

import (
	"context"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
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

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db}

}

func (r *quizRepository) CreateQuiz(ctx context.Context, quiz *model.Quiz) error {
	res := r.db.WithContext(ctx).Create(quiz)
	return res.Error
}

func (r *quizRepository) UpdateQuiz(ctx context.Context, quiz *model.Quiz) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(quiz).Updates(quiz).Error; err != nil {
			return err
		}
		if err := tx.Model(quiz).Association("Notes").Replace(quiz.Notes); err != nil {
			return err
		}

		return nil
	})
}

func (r *quizRepository) DeleteQuiz(ctx context.Context, quizID uuid.UUID) error {
	result := r.db.WithContext(ctx).Select("Questions").Delete(&model.Quiz{ID: quizID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *quizRepository) FindbyID(ctx context.Context, quizID uuid.UUID) (*model.Quiz, error) {
	var quiz model.Quiz
	result := r.db.WithContext(ctx).Preload("Questions").Preload("Notes").First(&quiz, quizID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &quiz, nil
}

func (r *quizRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Quiz, error) {
	var quizzes []*model.Quiz
	result := r.db.WithContext(ctx).Preload("Questions").Preload("Notes").Where("user_id = ?", userID).Find(&quizzes)
	return quizzes, result.Error
}

func (r *quizRepository) AddQuestion(ctx context.Context, question *model.Question) error {
	result := r.db.WithContext(ctx).Create(question)
	return result.Error
}

func (r *quizRepository) UpdateQuestion(ctx context.Context, question *model.Question) error {
	result := r.db.WithContext(ctx).Model(question).Updates(question)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *quizRepository) DeleteQuestion(ctx context.Context, questionID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Question{}, questionID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
