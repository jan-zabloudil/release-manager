package repository

import "github.com/nedpals/supabase-go"

type Repository struct {
	User     *UserRepository
	Settings *SettingsRepository
}

func NewRepository(client *supabase.Client) *Repository {
	return &Repository{
		User:     &UserRepository{client},
		Settings: NewSettingsRepository(client),
	}
}
