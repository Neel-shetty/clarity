package repository

import (
	"context"

	"github.com/HarshithRajesh/clarity/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	CheckUserExist(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositroy{db}
}

func (r *userRepository) CheckUserExist(ctx context.Context, email string) (string, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	res := r.db.WithContext(ctx).Create(user)
	return res.Error
}
