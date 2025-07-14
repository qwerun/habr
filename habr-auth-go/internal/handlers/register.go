package handlers

import (
	"encoding/json"
	"errors"
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/models"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"log"
	"net/http"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Bad JSON: %v, Body: %v", err, r.Body)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err = req.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashed, err := auth.HashPassword(req.PasswordHash)
	if err != nil {
		log.Printf("register: Failed to hash password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	req.PasswordHash = hashed

	user := models.NewUser(req.Email, req.PasswordHash, req.Nickname)
	var id string
	if id, err = s.explorer.Create(r.Context(), user); err != nil {
		switch {
		case errors.Is(err, user_repository.ErrEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, user_repository.ErrNicknameAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	code, err := s.explorer.SetVerificationCode(r.Context(), req.Email)
	if err != nil {
		log.Printf("Failed to SetVerificationCode: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = s.explorer.SendVerificationCode(req.Email, code)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	response := map[string]string{
		"id": id,
	}

	if err = writeJSON(w, response, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) verify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Bad JSON: %v, Body: %v", err, r.Body)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = s.explorer.CheckVerificationCode(r.Context(), req.Email, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrCodeCheckNotFound):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	err = s.explorer.VerifiedAccount(r.Context(), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrVerifyAccount):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	response := map[string]bool{
		"success": true,
	}
	if err = writeJSON(w, response, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
