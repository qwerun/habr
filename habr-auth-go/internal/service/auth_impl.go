package service

import (
	"context"
	"errors"
	"github.com/qwerun/habr-auth-go/internal/dto"
)

func (s *authService) Login(ctx context.Context, req dto.Login) (accessToken, refreshToken string, err error) {
	// TODO: implement Login logic
	return "", "", errors.New("Login not implemented")
}

func (s *authService) Refresh(ctx context.Context, req dto.Refresh) (newAccessToken, newRefreshToken string, err error) {
	// TODO: implement Refresh logic
	return "", "", errors.New("Refresh not implemented")
}
