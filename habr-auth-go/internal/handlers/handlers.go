package handlers

import (
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
	"github.com/qwerun/habr-auth-go/internal/service"
	"net/http"
)

//type Server struct {
//	explorer *user_repository.Repository
//	jwt      *auth.JwtManager
//}

//func NewMux(explorer *user_repository.Repository, jwt *auth.JwtManager) (http.Handler, error) {
//	server := &Server{explorer: explorer, jwt: jwt}
//	mux := http.NewServeMux()
//	mux.HandleFunc("/api/v1/register", server.register)
//	mux.HandleFunc("/api/v1/verify-email", server.verify)
//	mux.HandleFunc("/api/v1/login", server.login)
//	mux.HandleFunc("/api/v1/refresh", server.refresh)
//	mux.HandleFunc("/api/v1/change-password", server.changePassword)
//	return onlyPOST(mux), nil
//}

type Server struct {
	svc service.CompositeService
}

func NewMux(
	repo *user_repository.Repository,
	jwt *auth.JwtManager,
) (http.Handler, error) {

	server := &Server{
		svc: service.NewCompositeService(repo, jwt),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", server.register)
	mux.HandleFunc("/api/v1/verify-email", server.verify)
	mux.HandleFunc("/api/v1/login", server.login)
	mux.HandleFunc("/api/v1/refresh", server.refresh)
	mux.HandleFunc("/api/v1/change-password", server.changePassword)

	return onlyPOST(mux), nil
}
