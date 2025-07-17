package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"log"
	"net/http"
)

func (s *authService) Login(ctx context.Context, req dto.Login) (accessToken, refreshToken string, status int, err error) {
	hash, err := s.repo.GetPassHash(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrBadRequest):
			return "", "", http.StatusConflict, err
		default:
			return "", "", http.StatusInternalServerError, errors.New("Internal Server Error")
		}
	}

	err = auth.CheckPasswordHash(hash, req.Password)
	if err != nil {
		log.Printf("login: Failed to CheckPasswordHash: %v", err)
		return "", "", http.StatusBadRequest, errors.New("Bad JSON")
		return
	}

	id, err := s.repo.GetUserId(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrBadRequest):
			return "", "", http.StatusConflict, err
		default:
			log.Printf("login: Failed to GetUserId: %v", err)
			return "", "", http.StatusInternalServerError, errors.New("Internal Server Error")
		}
	}

	access, refresh, err := s.jwt.NewPair(id, req.FingerPrint)
	if err != nil {
		log.Printf("NewPair: Failed to NewPair: %v", err)
		return "", "", http.StatusInternalServerError, errors.New("Internal Server Error")
	}
	err = s.repo.SaveToken(ctx, id, req.FingerPrint, refresh, s.jwt.RefreshTtl)
	if err != nil {
		log.Printf("NewPair: Failed to SaveToken: %v", err)
		return "", "", http.StatusInternalServerError, errors.New("Internal Server Error")
	}
	return access, refresh, http.StatusOK, nil
}

func (s *authService) Refresh(ctx context.Context, req dto.Refresh) (newAccessToken, newRefreshToken string, status int, err error) {

	claims, err := s.jwt.ParseAccess(req.Access)
	switch {
	case err == nil:
	case errors.Is(err, jwt.ErrTokenExpired):
	default:
		log.Printf("refresh: ParseAccess error %v", err)
		return "", "", http.StatusUnauthorized, errors.New("Unauthorized")
	}

	userId := claims.Subject
	token, err := s.repo.GetToken(ctx, userId, req.FingerPrint)
	if err != nil || token != req.Refresh {
		log.Printf("refresh: Failed to GetToken: %v", err)
		return "", "", http.StatusUnauthorized, errors.New("Unauthorized")

	}
	access, refresh, err := s.jwt.NewPair(userId, req.FingerPrint)
	if err != nil {
		log.Printf("NewPair: Failed to NewPair: %v", err)
		return "", "", http.StatusUnauthorized, errors.New("Unauthorized")

	}
	err = s.repo.SaveToken(ctx, userId, req.FingerPrint, refresh, s.jwt.RefreshTtl)
	if err != nil {
		log.Printf("NewPair: Failed to SaveToken: %v", err)
		return "", "", http.StatusUnauthorized, errors.New("Unauthorized")
	}
	return access, refresh, http.StatusOK, nil
}
