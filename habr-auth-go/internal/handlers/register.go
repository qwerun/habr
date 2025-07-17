package handlers

import (
	"encoding/json"
	"github.com/qwerun/habr-auth-go/internal/dto"
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
	id, status, err := s.svc.Register(r.Context(), req)
	if err != nil {
		log.Printf("ERROR s.svc.Register: %v, Body: %v", err, r.Body)
		http.Error(w, err.Error(), status)
		return
	}

	response := map[string]string{
		"id": id,
	}

	if err = writeJSON(w, response, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	status, err := s.svc.Verify(r.Context(), req.Email, req.Code)
	if err != nil {
		log.Printf("ERROR s.svc.Verify: %v, Body: %v", err, r.Body)
		http.Error(w, err.Error(), status)
		return
	}

	response := map[string]bool{
		"success": true,
	}
	if err = writeJSON(w, response, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
