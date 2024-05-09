package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"
)

type Repository struct {
	Auth     *AuthRepository
	User     *UserRepository
	Project  *ProjectRepository
	Settings *SettingsRepository
	Release  *ReleaseRepository
}

func NewRepository(client *supabase.Client, pool *pgxpool.Pool) *Repository {
	return &Repository{
		Auth:     NewAuthRepository(client),
		User:     NewUserRepository(client),
		Project:  NewProjectRepository(client),
		Settings: NewSettingsRepository(client),
		Release:  NewReleaseRepository(pool),
	}
}
