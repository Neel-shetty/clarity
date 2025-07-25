package domain

import "time"

type User struct {
	Id            uuid   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	Password_hash string `gorm:"not null"`
	Created_at    time.Time
}
