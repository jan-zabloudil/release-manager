package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Release struct {
	ID           uuid.UUID `json:"id"`
	ProjectID    uuid.UUID `json:"project_id"`
	ReleaseTitle string    `json:"release_title"`
	ReleaseNotes string    `json:"release_notes"`
	AuthorUserID uuid.UUID `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToRelease(r svcmodel.Release) Release {
	return Release{
		ID:           r.ID,
		ProjectID:    r.ProjectID,
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
		AuthorUserID: r.AuthorUserID,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
