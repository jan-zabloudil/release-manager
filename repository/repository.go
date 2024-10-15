package repository

import (
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"
)

type githubURLGenerator interface {
	GenerateRepoURL(ownerSlug, repoSlug string) (url.URL, error)
	GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error)
}

type fileURLGenerator interface {
	GenerateFileURL(filePath string) (url.URL, error)
}

type Repository struct {
	User     *UserRepository
	Project  *ProjectRepository
	Settings *SettingsRepository
	Release  *ReleaseRepository
}

func NewRepository(
	client *supabase.Client,
	pool *pgxpool.Pool,
	githubURLGenerator githubURLGenerator,
	fileURLGenerator fileURLGenerator,
) *Repository {
	return &Repository{
		User:     NewUserRepository(client, pool),
		Project:  NewProjectRepository(pool, githubURLGenerator),
		Settings: NewSettingsRepository(pool),
		Release:  NewReleaseRepository(pool, githubURLGenerator, fileURLGenerator),
	}
}
