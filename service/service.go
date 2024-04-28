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
	gc model.GithubClient,
) *Service {
	authSvc := NewAuthService(ar, ur)
	settingsSvc := NewSettingsService(authSvc, sr)
	projectSvc := NewProjectService(authSvc, settingsSvc, pr, env, pi, gc)

	return &Service{
		Auth:              authSvc,
		User:              NewUserService(authSvc, ur),
		Project:           projectSvc,
		Settings:          settingsSvc,
		ProjectMembership: NewProjectMembershipService(authSvc, projectSvc, pi),
	}
}
