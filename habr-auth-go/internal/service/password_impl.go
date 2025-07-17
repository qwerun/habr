package service

import (
	"context"
	"errors"
	"github.com/qwerun/habr-auth-go/internal/dto"
)

func (s *passwordService) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error {
	// TODO: implement ChangePassword logic
	return errors.New("ChangePassword not implemented")
}
