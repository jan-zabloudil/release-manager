package service

import "release-manager/service/model"

type Service struct {
	Auth *AuthService
	User *UserService
}

func NewService(ar model.AuthRepository, ur model.UserRepository) *Service {
	authSvc := NewAuthService(ar, ur)

	return &Service{
		Auth: authSvc,
		User: NewUserService(authSvc, ur),
	}
}
