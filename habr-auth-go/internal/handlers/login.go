package handlers

import (
	"encoding/json"
	"github.com/qwerun/habr-auth-go/internal/dto"
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
	access, refresh, status, err := s.svc.Login(r.Context(), req)
	if err != nil {
		log.Printf("Login: s.svc.Login: error %v", err)
		http.Error(w, err.Error(), status)
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
	}
	if err = writeJSON(w, response, status); err != nil {
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
	access, refresh, status, err := s.svc.Refresh(r.Context(), req)
	if err != nil {
		log.Printf("Refresh: s.svc.Refresh: error %v", err)
		http.Error(w, err.Error(), status)
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
	}
	if err = writeJSON(w, response, status); err != nil {
		log.Printf("refresh: final error %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}
