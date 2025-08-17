package service

import (
	"context"
	"errors"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/Neel-shetty/clarity/internal/repository"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Signup(ctx context.Context, signup *model.SignUp) error
	Login(ctx context.Context, login *model.Login) (*model.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*model.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Signup(ctx context.Context, signup *model.SignUp) error {
	if signup.Password != signup.CheckPassword {
		return errors.New("passwords do not match")
	}

	exists, _ := s.repo.CheckUserExist(ctx, signup.Email)
	if exists != nil {
		return errors.New("user with this email already exists")
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
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

func (s *userService) Login(ctx context.Context, login *model.Login) (*model.User, error) {
	user, err := s.repo.CheckUserExist(ctx, login.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email ")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		return nil, errors.New("invalid  password")
	}
	return user, nil
}
func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return user, nil
}
