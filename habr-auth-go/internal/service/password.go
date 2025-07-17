package service

import (
	"context"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"log"
)

type PasswordService interface {
	ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) (int, error)
}

type passwordService struct {
	repo *user_repository.Repository
}

func NewPasswordService(repo *user_repository.Repository) PasswordService {
	if repo == nil {
		log.Fatal("NewPasswordService: repository cannot be nil")
	}
	return &passwordService{
		repo: repo,
	}
}
