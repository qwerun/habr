package service

import (
	"context"
	"github.com/qwerun/habr-auth-go/internal/dto"
)

type RegistrationService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (userID string, err error)
	Verify(ctx context.Context, email, code string) error
}

type registrationService struct {
}

func NewRegistrationService() RegistrationService {
	return &registrationService{}
}
