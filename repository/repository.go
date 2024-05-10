package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"
)

type Repository struct {
	User     *UserRepository
	Project  *ProjectRepository
	Settings *SettingsRepository
	Release  *ReleaseRepository
}

func NewRepository(client *supabase.Client, pool *pgxpool.Pool) *Repository {
	return &Repository{
		User:     NewUserRepository(client),
		Project:  NewProjectRepository(client),
		Settings: NewSettingsRepository(client),
		Release:  NewReleaseRepository(pool),
	}
}
