package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/datatypes"
)

type Note struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	FileName  string    `gorm:"type:text;not null"`
	S3Key     string    `gorm:"type:text;uniqueIndex;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
}

type Quiz struct {
	ID        uuid.UUID  `gorm:"type:uuid; primaryKey;"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	Title     string     `gorm:"type:text;"`
	Notes     []*Note    `gorm:"many2many:quiz_notes;"`
	User      User       `gorm:"foreignKey:UserID"`
	Questions []Question `gorm:"foreignKey:QuizID"`
}

type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Question struct {
	ID            uuid.UUID                   `gorm:"type:uuid;primaryKey;"`
	QuizID        uuid.UUID                   `gorm:"type:uuid;not null"`
	QuestionText  string                      `gorm:"type:text;not null"`
	QuestionType  string                      `gorm:"type:varchar(50);not null"`
	Points        int                         `gorm:"not null"`
	Options       datatypes.JSONSlice[Option] `gorm:"type:jsonb;not null"`
	CorrectOption string                      `gorm:"type:text;not null"`
}

type QuizAttempt struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	QuizID      uuid.UUID `gorm:"type:uuid;not null"`
	Score       int
	StartedAt   time.Time
	CompletedAt *time.Time
	User        User `gorm:"foreignKey:UserID"`
	Quiz        Quiz `gorm:"foreignKey:QuizID"`
}
