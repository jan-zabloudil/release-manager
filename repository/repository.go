package repository

import "github.com/nedpals/supabase-go"

type Repository struct {
	Auth     *AuthRepository
	User     *UserRepository
	Project  *ProjectRepository
	Settings *SettingsRepository
}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{
		Auth:     NewAuthRepository(client),
		User:     NewUserRepository(client),
		Project:  NewProjectRepository(client),
		Settings: NewSettingsRepository(client),
	}
}
