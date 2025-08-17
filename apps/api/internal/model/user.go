package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Profile struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	FullName  string    `gorm:"type:varchar(255)"`
	AvatarURL string    `gorm:"type:text"`
	Bio       string    `gorm:"type:text"`
}

type SignUp struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	CheckPassword string `json:"check_password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
