package service

import (
	"context"
	"errors"

	"github.com/HarshitRajesh/clarity/internal/domain"
	"github.com/HarshitRajesh/clarity/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Signup(signup *domain.Signup) error
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Signup(ctx context.Context, signup *domain.Signup) error {
	if signup.Password != signup.CheckPassword {
		return nil, errors.New("passwords do not match")
	}

	email := s.repo.CheckUserExist(ctx, signup.email)
	hashPass, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Name:     signup.Name,
		Email:    signup.Email,
		Password: hashPass,
	}

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
