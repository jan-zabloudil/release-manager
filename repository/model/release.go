package model

import (
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
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
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
