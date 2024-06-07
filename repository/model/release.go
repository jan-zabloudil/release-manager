package model

import (
	"database/sql"
	"fmt"
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

type tagURLGeneratorFunc func(ownerSlug, repoSlug, tag string) (url.URL, error)

func ToSvcRelease(rls Release, urlGenerator tagURLGeneratorFunc) (svcmodel.Release, error) {
	tagURL, err := urlGenerator(rls.GithubOwnerSlug.String, rls.GithubRepoSlug.String, rls.GitTagName)
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("failed to generate tag URL: %w", err)
	}

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
	}, nil
}

func ToSvcReleases(releases []Release, urlGenerator tagURLGeneratorFunc) ([]svcmodel.Release, error) {
	r := make([]svcmodel.Release, 0, len(releases))
	for _, release := range releases {
		svcRls, err := ToSvcRelease(release, urlGenerator)
		if err != nil {
			return nil, err
		}

		r = append(r, svcRls)
	}

	return r, nil
}
