package service

import (
	"context"
	"errors"
	"github.com/qwerun/habr-auth-go/internal/dto"
)

func (s *registrationService) Register(ctx context.Context, req dto.RegisterRequest) (userID string, err error) {
	// TODO: implement Register logic
	return "", errors.New("Register not implemented")
}

func (s *registrationService) Verify(ctx context.Context, email, code string) error {
	// TODO: implement Verify logic
	return errors.New("Verify not implemented")
}
