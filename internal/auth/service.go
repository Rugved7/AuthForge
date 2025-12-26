package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Rugved7/authforge/internal/user"
	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("Invalid credentials: check email or password")

type Service struct {
	users user.Repository
}

func NewService(users user.Repository) *Service {
	return &Service{
		users: users,
	}
}

func (s *Service) Signup(ctx context.Context, email, password string) (*user.User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         user.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*user.User, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := ComparePassword(password, u.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}
