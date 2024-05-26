package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Release struct {
	ID                uuid.UUID         `db:"id"`
	ProjectID         uuid.UUID         `db:"project_id"`
	ReleaseTitle      string            `db:"release_title"`
	ReleaseNotes      string            `db:"release_notes"`
	AuthorUserID      uuid.UUID         `db:"created_by"`
	CreatedAt         time.Time         `db:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at"`
	GithubReleaseID   int64             `db:"github_release_id"`
	GithubOwnerSlug   string            `db:"github_owner_slug"`
	GithubRepoSlug    string            `db:"github_repo_slug"`
	GithubReleaseData GithubReleaseData `db:"github_release_data"`
}

type GithubReleaseData struct {
	GitTagName string    `json:"git_tag_name"`
	HTMLURL    string    `json:"html_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToSvcRelease(rls Release) svcmodel.Release {
	return svcmodel.Release{
		ID:           rls.ID,
		ProjectID:    rls.ProjectID,
		ReleaseTitle: rls.ReleaseTitle,
		ReleaseNotes: rls.ReleaseNotes,
		AuthorUserID: rls.AuthorUserID,
		CreatedAt:    rls.CreatedAt,
		UpdatedAt:    rls.UpdatedAt,
	}
}

func ToSvcReleases(releases []Release) []svcmodel.Release {
	r := make([]svcmodel.Release, 0, len(releases))
	for _, release := range releases {
		r = append(r, ToSvcRelease(release))
	}
	return r
}
