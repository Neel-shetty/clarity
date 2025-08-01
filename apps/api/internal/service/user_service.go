package service

import (
	"context"
	"errors"

	"github.com/HarshitRajesh/clarity/internal/domain"
	"github.com/HarshitRajesh/clarity/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Signup(signup *domain.Signup) error
	Login(ctx context.Context, login *domain.Login) (*domain.User, error)
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

	_, err := s.repo.CheckUserExist(ctx, signup.email)
	if err != nil {
		return nil, errors.New("User with this email already exists")
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Name:     signup.Name,
		Email:    signup.Email,
		Password: string(hashPass),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(ctx context.Context, login *domain.Login) (*domain.User, error) {
	user, err := s.repo.CheckUserExist(ctx, login.email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword(login.Password, password)
	if err != nil {
		return errors.New("invalid email or password")
	}
	return user, nil
}
