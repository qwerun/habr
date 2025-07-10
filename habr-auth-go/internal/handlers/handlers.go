package handlers

import (
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"net/http"
)

type Server struct {
	explorer *user_repository.Repository
}

func NewMux(explorer *user_repository.Repository) (http.Handler, error) {
	server := &Server{explorer: explorer}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", server.register)
	mux.HandleFunc("/api/v1/verify-email", server.verify)
	mux.HandleFunc("/api/v1/login", server.defaults)
	mux.HandleFunc("/api/v1/change-password", server.changePassword)
	return onlyPOST(mux), nil
}
