package service

import (
	"context"
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
)

type AuthService interface {
	Login(ctx context.Context, req dto.Login) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, req dto.Refresh) (newAccessToken, newRefreshToken string, err error)
}

type authService struct {
	repo *user_repository.Repository
	jwt  *auth.JwtManager
}

func NewAuthService(repo *user_repository.Repository, jwt *auth.JwtManager) AuthService {
	return &authService{repo: repo, jwt: jwt}
}
