package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"
)

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type Repository struct {
	User     *UserRepository
	Project  *ProjectRepository
	Settings *SettingsRepository
	Release  *ReleaseRepository
}

func NewRepository(client *supabase.Client, pool *pgxpool.Pool) *Repository {
	return &Repository{
		User:     NewUserRepository(client, pool),
		Project:  NewProjectRepository(pool),
		Settings: NewSettingsRepository(pool),
		Release:  NewReleaseRepository(pool),
	}
}
