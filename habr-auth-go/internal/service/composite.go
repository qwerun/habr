package service

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

func NewCompositeService() CompositeService {
	return &compositeService{
		PasswordService:     NewPasswordService(),
		AuthService:         NewAuthService(),
		RegistrationService: NewRegistrationService(),
	}
}
