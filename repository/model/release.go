package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Release struct {
	ID                     uuid.UUID `db:"id"`
	ProjectID              uuid.UUID `db:"project_id"`
	ReleaseTitle           string    `db:"release_title"`
	ReleaseNotes           string    `db:"release_notes"`
	AuthorUserID           uuid.UUID `db:"created_by"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
	GitTagName             string    `db:"git_tag_name"`
	GithubReleaseCreatedAt time.Time `db:"github_release_created_at"`
	GithubReleaseUpdatedAt time.Time `db:"github_release_updated_at"`
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
