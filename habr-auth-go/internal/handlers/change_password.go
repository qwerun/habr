package handlers

import (
	"encoding/json"
	"errors"
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/dto"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"log"
	"net/http"
)

func (s *Server) changePassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ChangePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Bad JSON: %v, Body: %v", err, r.Body)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}
	if err = req.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hash, err := s.explorer.GetPassHash(req.Email)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = auth.CheckPasswordHash(hash, req.Password)
	if err != nil {
		log.Printf("changePassword: Failed to CheckPasswordHash: %v", err)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}
	newHashed, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("changePassword: Failed to hash password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = s.explorer.SetNewHash(req.Email, newHashed)
	if err != nil {
		log.Printf("SetNewHash: Failed to set new hash password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := map[string]bool{
		"success": true,
	}
	if err = writeJSON(w, response, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
