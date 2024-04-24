package service

import "release-manager/service/model"

type Service struct {
	Auth              *AuthService
	User              *UserService
	Project           *ProjectService
	Settings          *SettingsService
	ProjectMembership *ProjectMembershipService
}

func NewService(
	ar model.AuthRepository,
	ur model.UserRepository,
	pr model.ProjectRepository,
	env model.EnvironmentRepository,
	sr model.SettingsRepository,
	pi model.ProjectInvitationRepository,
) *Service {
	authSvc := NewAuthService(ar, ur)
	projectSvc := NewProjectService(authSvc, pr, env, pi)

	return &Service{
		Auth:              authSvc,
		User:              NewUserService(authSvc, ur),
		Project:           projectSvc,
		Settings:          NewSettingsService(authSvc, sr),
		ProjectMembership: NewProjectMembershipService(authSvc, projectSvc, pi),
	}
}
