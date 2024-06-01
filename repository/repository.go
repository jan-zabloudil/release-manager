package repository

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"
)

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type githubURLGenerator interface {
	GenerateRepoURL(ownerSlug, repoSlug string) (url.URL, error)
	GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error)
}

type Repository struct {
	User     *UserRepository
	Project  *ProjectRepository
	Settings *SettingsRepository
	Release  *ReleaseRepository
}

func NewRepository(client *supabase.Client, pool *pgxpool.Pool, urlGenerator githubURLGenerator) *Repository {
	return &Repository{
		User:     NewUserRepository(client, pool),
		Project:  NewProjectRepository(pool, urlGenerator),
		Settings: NewSettingsRepository(pool),
		Release:  NewReleaseRepository(pool, urlGenerator),
	}
}
