package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/Neel-shetty/clarity/internal/repository"
	"github.com/Neel-shetty/clarity/internal/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofrs/uuid/v5"
)

type NoteService interface {
	CreateNote(ctx context.Context, userID uuid.UUID, fileName string, contentType string) (*model.Note, *s3.PresignedPostRequest, error)
	GetNoteByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type noteService struct {
	repo    repository.NoteRepository
	storage storage.Storage
}

func NewNoteService(repo repository.NoteRepository, storage storage.Storage) NoteService {
	return &noteService{repo, storage}
}

var allowedTypes = []string{
	// PDF
	"application/pdf",

	// Microsoft Word
	"application/msword", // .doc
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document", // .docx

	// Microsoft PowerPoint
	"application/vnd.ms-powerpoint",                                             // .ppt
	"application/vnd.openxmlformats-officedocument.presentationml.presentation", // .pptx
}

var mb10 = 10 * 1024 * 1024

func (s *noteService) CreateNote(ctx context.Context, userID uuid.UUID, fileName string, contentType string) (*model.Note, *s3.PresignedPostRequest, error) {
	key := fmt.Sprintf("%s%s", userID.String(), fileName)
	note := &model.Note{
		UserID:   userID,
		S3Key:    key,
		FileName: fileName,
	}

	if err := s.repo.CreateNote(ctx, note); err != nil {
		return nil, &s3.PresignedPostRequest{}, err
	}

	presignedReq, err := s.storage.GenerateUploadURL(ctx, key, contentType, allowedTypes, mb10, 10*time.Minute)
	if err != nil {
		return nil, &s3.PresignedPostRequest{}, err
	}

	return note, presignedReq, nil
}

func (s *noteService) GetNoteByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Note, error) {
	note, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if note.UserID != userID {
		return nil, fmt.Errorf("unauthorized: you do not own this note")
	}
	return note, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	note, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if note.UserID != userID {
		return fmt.Errorf("unauthorized: you do not own this note")
	}

	if err := s.repo.DeleteNote(ctx, id); err != nil {
		return err
	}
	return nil
}
