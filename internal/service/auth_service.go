package service

import (
	"backend-store/internal/auth"
	"backend-store/internal/repository"
	"context"
	"fmt"
)

type AuthService interface {
	Register(ctx context.Context, email, password, firstName, lastName string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type authService struct {
	authClient *auth.AuthClient
	userRepo   repository.UserRepository
}

func NewAuthService(authClient *auth.AuthClient, userRepo repository.UserRepository) AuthService {
	return &authService{
		authClient: authClient,
		userRepo:   userRepo,
	}
}

func (s *authService) Register(ctx context.Context, email, password, firstName, lastName string) (string, error) {
	registerResp, err := s.authClient.Register(ctx, email, password, firstName, lastName)
	if err != nil {
		return "", err
	}

	// Временно возвращаем успешную регистрацию без токена
	// TODO: Реализовать метод Login в auth-service для получения токена после регистрации
	return fmt.Sprintf("registration_success_%d", registerResp.Id), nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	resp, err := s.authClient.Login(ctx, email, password)
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (string, error) {
	resp, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return "", err
	}

	if !resp.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return resp.UserId, nil
}
