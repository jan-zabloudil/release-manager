package repository

import "github.com/nedpals/supabase-go"

type Repository struct {
	Auth              *AuthRepository
	User              *UserRepository
	Project           *ProjectRepository
	Environment       *EnvironmentRepository
	Settings          *SettingsRepository
	ProjectInvitation *ProjectInvitationRepository
}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{
		Auth:              NewAuthRepository(client),
		User:              NewUserRepository(client),
		Project:           NewProjectRepository(client),
		Environment:       NewEnvironmentRepository(client),
		Settings:          NewSettingsRepository(client),
		ProjectInvitation: NewProjectInvitationRepository(client),
	}
}
