package service

import (
	"context"
	"errors"
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/models"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"log"
	"net/http"
)

func (s *registrationService) Register(ctx context.Context, req dto.RegisterRequest) (userID string, status int, err error) {
	hashed, err := auth.HashPassword(req.PasswordHash)
	if err != nil {
		log.Printf("register: Failed to hash password: %v", err)
		return "", http.StatusInternalServerError, errors.New("Internal Server Error")
	}
	req.PasswordHash = hashed

	user := models.NewUser(req.Email, req.PasswordHash, req.Nickname)
	var id string
	if id, err = s.repo.Create(ctx, user); err != nil {
		switch {
		case errors.Is(err, user_repository.ErrEmailAlreadyExists):
			return "", http.StatusConflict, err
		case errors.Is(err, user_repository.ErrNicknameAlreadyExists):
			return "", http.StatusConflict, err
		default:
			return "", http.StatusInternalServerError, errors.New("Internal Server Error")
		}
	}

	code, err := s.repo.SetVerificationCode(ctx, req.Email)
	if err != nil {
		log.Printf("Failed to SetVerificationCode: %v", err)
		return "", http.StatusInternalServerError, errors.New("Internal Server Error")
	}

	err = s.repo.SendVerificationCode(req.Email, code)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("Internal Server Error")
	}
	return id, http.StatusOK, errors.New("Register not implemented")
}

func (s *registrationService) Verify(ctx context.Context, email, code string) (int, error) {

	err := s.repo.CheckVerificationCode(ctx, email, code)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrCodeCheckNotFound):
			return http.StatusConflict, err
		default:
			log.Printf("Failed to CheckVerificationCode: %v", err)
			return http.StatusInternalServerError, errors.New("Internal Server Error")
		}
	}
	err = s.repo.VerifiedAccount(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrVerifyAccount):
			return http.StatusConflict, err
		default:
			log.Printf("Failed to VerifiedAccount: %v", err)
			return http.StatusInternalServerError, errors.New("Internal Server Error")
		}
	}
	return http.StatusOK, nil
}
