package repository

import (
	"context"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	CheckUserExist(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CheckUserExist(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	res := r.db.WithContext(ctx).Create(user)
	return res.Error
}

func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id=?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
