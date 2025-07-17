package service

import (
	"context"
	"github.com/qwerun/habr-auth-go/internal/dto"
)

type PasswordService interface {
	ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error
}

type passwordService struct {
}

func NewPasswordService() PasswordService {
	return &passwordService{}
}
