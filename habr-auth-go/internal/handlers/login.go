package handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
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
	hash, err := s.explorer.GetPassHash(r.Context(), req.Email)
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

	id, err := s.explorer.GetUserId(r.Context(), req.Email)
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
	err = s.explorer.SaveToken(r.Context(), id, req.FingerPrint, refresh, s.jwt.RefreshTtl)
	if err != nil {
		log.Printf("NewPair: Failed to SaveToken: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
		FP      string `json:"fingerprint"`
	}{
		Access:  access,
		Refresh: refresh,
		FP:      req.FingerPrint,
	} // если бы это был не тестовый проект, то передавал бы refresh в куках с пометкой HttpOnly true и Secure true
	if err = writeJSON(w, response, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.Refresh
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Bad JSON: %v, Body: %v", err, r.Body)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()
	if err = req.IsValid(); err != nil {
		log.Printf("refresh: is valid error %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := s.jwt.ParseAccess(req.Access)
	switch {
	case err == nil:
	case errors.Is(err, jwt.ErrTokenExpired):
	default:
		log.Printf("refresh: ParseAccess error %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	userId := claims.Subject
	token, err := s.explorer.GetToken(r.Context(), userId, req.FingerPrint)
	if err != nil || token != req.Refresh {
		log.Printf("refresh: Failed to GetToken: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	access, refresh, err := s.jwt.NewPair(userId, req.FingerPrint)
	if err != nil {
		log.Printf("NewPair: Failed to NewPair: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err = s.explorer.SaveToken(r.Context(), userId, req.FingerPrint, refresh, s.jwt.RefreshTtl)
	if err != nil {
		log.Printf("NewPair: Failed to SaveToken: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	response := struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
		FP      string `json:"fingerprint"`
	}{
		Access:  access,
		Refresh: refresh,
		FP:      req.FingerPrint,
	} // если бы это был не тестовый проект, то передавал бы refresh в куках с пометкой HttpOnly true и Secure true
	if err = writeJSON(w, response, http.StatusOK); err != nil {
		log.Printf("refresh: final error %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
