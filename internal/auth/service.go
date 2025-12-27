package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Rugved7/authforge/internal/user"
	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid credentials: check email or password")

type Service struct {
	users        user.Repository
	tokenManager *TokenManager
}

type AuthResult struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
}

func NewService(users user.Repository, tokenManager *TokenManager) *Service {
	return &Service{
		users:        users,
		tokenManager: tokenManager,
	}
}

func (s *Service) Signup(ctx context.Context, email, password string) (*AuthResult, error) {
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

	accessToken, err := s.tokenManager.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := ComparePassword(password, u.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	_, claims, err := s.tokenManager.ParseToken(refreshToken)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Ensure refresh token
	if t, ok := claims["type"].(string); !ok || t != "refresh" {
		return "", ErrInvalidCredentials
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", ErrInvalidCredentials
	}

	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return s.tokenManager.GenerateAccessToken(u.ID, string(u.Role))
}
