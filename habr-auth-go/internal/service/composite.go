package service

import (
	"github.com/qwerun/habr-auth-go/internal/auth"
	"github.com/qwerun/habr-auth-go/internal/repository/user_repository"
)

type CompositeService interface {
	PasswordService
	AuthService
	RegistrationService
}

type compositeService struct {
	PasswordService
	AuthService
	RegistrationService
}

//func NewCompositeService() CompositeService {
//	return &compositeService{
//		PasswordService:     NewPasswordService(),
//		AuthService:         NewAuthService(),
//		RegistrationService: NewRegistrationService(),
//	}
//}

func NewCompositeService(
	repo *user_repository.Repository,
	jwt *auth.JwtManager,
) CompositeService {
	return &compositeService{
		PasswordService:     NewPasswordService(repo),
		AuthService:         NewAuthService(repo, jwt),
		RegistrationService: NewRegistrationService(repo),
	}
}
