package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateReleaseInput struct {
	ReleaseTitle string `json:"release_title"`
	ReleaseNotes string `json:"release_notes"`
	GitTagName   string `json:"git_tag_name"`
}

type UpdateReleaseInput struct {
	ReleaseTitle *string `json:"release_title"`
	ReleaseNotes *string `json:"release_notes"`
}

type Release struct {
	ID           uuid.UUID `json:"id"`
	ReleaseTitle string    `json:"release_title"`
	ReleaseNotes string    `json:"release_notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToSvcCreateReleaseInput(r CreateReleaseInput) svcmodel.CreateReleaseInput {
	return svcmodel.CreateReleaseInput{
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
		GitTagName:   r.GitTagName,
	}
}

func ToSvcUpdateReleaseInput(r UpdateReleaseInput) svcmodel.UpdateReleaseInput {
	return svcmodel.UpdateReleaseInput{
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
	}
}

func ToRelease(r svcmodel.Release) Release {
	return Release{
		ID:           r.ID,
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func ToReleases(releases []svcmodel.Release) []Release {
	r := make([]Release, 0, len(releases))
	for _, release := range releases {
		r = append(r, ToRelease(release))
	}
	return r
}
