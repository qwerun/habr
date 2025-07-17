package service

import (
	"context"
	"errors"
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"log"
	"net/http"
)

func (s *passwordService) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) (int, error) {
	hash, err := s.repo.GetPassHash(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrBadRequest):
			return http.StatusConflict, err
		default:
			return http.StatusInternalServerError, errors.New("Internal Server Error")
		}
	}

	err = auth.CheckPasswordHash(hash, req.Password)
	if err != nil {
		log.Printf("changePassword: Failed to CheckPasswordHash: %v", err)
		return http.StatusBadRequest, errors.New("Bad JSON")
	}
	newHashed, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("changePassword: Failed to hash password: %v", err)
		return http.StatusInternalServerError, errors.New("Internal Server Error")
	}
	err = s.repo.SetNewHash(ctx, req.Email, newHashed)
	if err != nil {
		log.Printf("SetNewHash: Failed to set new hash password: %v", err)
		return http.StatusInternalServerError, errors.New("Internal Server Error")
	}
	return http.StatusOK, nil
}
