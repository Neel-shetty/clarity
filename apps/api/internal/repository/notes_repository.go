package repository

import (
	"context"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/gofrs/uuid/v5"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note *model.Note) error
	FindByID(ctx context.Context, noteID uuid.UUID) (*model.Note, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Note, error)
	DeleteNote(ctx context.Context, noteID uuid.UUID) error
}
