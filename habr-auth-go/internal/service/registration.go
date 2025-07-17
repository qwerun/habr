package service

import (
	"context"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
)

type RegistrationService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (userID string, status int, err error)
	Verify(ctx context.Context, email, code string) (int, error)
}

type registrationService struct {
	repo *user_repository.Repository
}

func NewRegistrationService(repo *user_repository.Repository) RegistrationService {
	return &registrationService{repo: repo}
}
