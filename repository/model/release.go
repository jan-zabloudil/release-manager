package model

import (
	"database/sql"
	"net/url"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Release struct {
	ID           uuid.UUID `db:"id"`
	ProjectID    uuid.UUID `db:"project_id"`
	ReleaseTitle string    `db:"release_title"`
	ReleaseNotes string    `db:"release_notes"`
	AuthorUserID uuid.UUID `db:"created_by"`
	GitTagName   string    `db:"git_tag_name"`
	// GithubRepoSlug and GithubOwnerSlug are fetched from the project
	// and are used to generate the tag URL
	GithubRepoSlug  sql.NullString `db:"github_repo_slug"`
	GithubOwnerSlug sql.NullString `db:"github_owner_slug"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

func ToSvcRelease(rls Release, tagURL url.URL) svcmodel.Release {
	return svcmodel.Release{
		ID:           rls.ID,
		ProjectID:    rls.ProjectID,
		ReleaseTitle: rls.ReleaseTitle,
		ReleaseNotes: rls.ReleaseNotes,
		GitTagName:   rls.GitTagName,
		GitTagURL:    tagURL,
		AuthorUserID: rls.AuthorUserID,
		CreatedAt:    rls.CreatedAt,
		UpdatedAt:    rls.UpdatedAt,
	}
}
