package service

import (
	"context"

	"github.com/qwerun/habr-auth-go/internal/dto"
)

type PasswordService interface {
	ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error
}

type AuthService interface {
	Login(ctx context.Context, req dto.Login) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, req dto.Refresh) (newAccessToken, newRefreshToken string, err error)
}

type RegistrationService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (userID string, err error)
	Verify(ctx context.Context, email, code string) error
}

type CompositeService interface {
	PasswordService
	AuthService
	RegistrationService
}
