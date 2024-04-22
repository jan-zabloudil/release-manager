package service

import "release-manager/service/model"

type Service struct {
	Auth    *AuthService
	User    *UserService
	Project *ProjectService
}

func NewService(ar model.AuthRepository, ur model.UserRepository, pr model.ProjectRepository, env model.EnvironmentRepository) *Service {
	authSvc := NewAuthService(ar, ur)

	return &Service{
		Auth:    authSvc,
		User:    NewUserService(authSvc, ur),
		Project: NewProjectService(authSvc, pr, env),
	}
}
