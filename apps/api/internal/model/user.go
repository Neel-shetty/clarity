package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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
