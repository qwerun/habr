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

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req dto.Login
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
		log.Printf("login: Failed to CheckPasswordHash: %v", err)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	id, err := s.explorer.GetUserId(req.Email)
	if err != nil {
		switch {
		case errors.Is(err, user_repository.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	access, refresh, err := s.jwt.NewPair(id, req.FingerPrint)
	if err != nil {
		log.Printf("NewPair: Failed to NewPair: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = s.explorer.SaveToken(id, req.FingerPrint, refresh, s.jwt.RefreshTtl)
	if err != nil {
		log.Printf("NewPair: Failed to SaveToken: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
		FP      string `json:"fp"`
	}{
		Access:  access,
		Refresh: refresh,
		FP:      req.FingerPrint,
	}
	if err = writeJSON(w, response, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
