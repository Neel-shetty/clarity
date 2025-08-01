package service

import (
	"context"
	"errors"

	"github.com/Neel-shetty/clarity/internal/domain"
	"github.com/Neel-shetty/clarity/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Signup(ctx context.Context, signup *domain.SignUp) error
	Login(ctx context.Context, login *domain.Login) error
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Signup(ctx context.Context, signup *domain.SignUp) error {
	if signup.Password != signup.CheckPassword {
		return errors.New("passwords do not match")
	}

	exists, _ := s.repo.CheckUserExist(ctx, signup.Email)
	if exists != nil {
		return errors.New("User with this email already exists")
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &domain.User{
		Name:         signup.Name,
		Email:        signup.Email,
		PasswordHash: string(hashPass),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(ctx context.Context, login *domain.Login) error {
	user, err := s.repo.CheckUserExist(ctx, login.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid email ")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		return errors.New("invalid  password")
	}
	return nil
}
