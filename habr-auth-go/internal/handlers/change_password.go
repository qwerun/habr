package handlers

import (
	"encoding/json"
	"github.com/qwerun/habr-auth-go/internal/dto"
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
	defer r.Body.Close()
	if err = req.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := s.svc.ChangePassword(r.Context(), req)
	if err != nil {
		log.Printf("ERROR s.svc.ChangePassword: %v, Body: %v", err, r.Body)
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
