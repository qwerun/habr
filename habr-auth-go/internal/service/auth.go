package service

import (
	"context"
	"github.com/qwerun/habr-auth-go/internal/dto"
)

type AuthService interface {
	Login(ctx context.Context, req dto.Login) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, req dto.Refresh) (newAccessToken, newRefreshToken string, err error)
}

type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}
