package repository

import (
	"context"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note *model.Note) error
	FindByID(ctx context.Context, noteID uuid.UUID) (*model.Note, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Note, error)
	DeleteNote(ctx context.Context, noteID uuid.UUID) error
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) CreateNote(ctx context.Context, note *model.Note) error {
	result := r.db.WithContext(ctx).Create(note)
	return result.Error
}

func (r *noteRepository) FindByID(ctx context.Context, noteID uuid.UUID) (*model.Note, error) {
	var note model.Note
	result := r.db.WithContext(ctx).Preload("User").First(&note, noteID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &note, nil
}

func (r *noteRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&notes)
	return notes, result.Error
}

func (r *noteRepository) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM quiz_notes WHERE note_id = ?", noteID).Error; err != nil {
			return err
		}

		result := tx.Delete(&model.Note{}, noteID)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}
