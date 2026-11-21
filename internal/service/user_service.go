package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"qubix/internal/model"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, user *model.User, password string) error {
	if strings.TrimSpace(user.Login) == "" || strings.TrimSpace(user.Email) == "" || password == "" {
		return errors.New("login, email and password are required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	return s.repo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, loginOrEmail, password string) (*model.User, error) {
	user, err := s.repo.GetByLoginOrEmail(ctx, loginOrEmail)
	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	if user.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if strings.TrimSpace(user.Login) == "" || strings.TrimSpace(user.Email) == "" {
		return errors.New("login and email are required")
	}

	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	return s.repo.Delete(ctx, id)
}
