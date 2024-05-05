package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	errReleaseTitleRequired = errors.New("release title is required")
)

type CreateReleaseInput struct {
	ReleaseTitle string
	ReleaseNotes string
}

type Release struct {
	ID           uuid.UUID
	ProjectID    uuid.UUID
	ReleaseTitle string
	ReleaseNotes string
	AuthorUserID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// TODO add source code link
}

func NewRelease(input CreateReleaseInput, projectID, authorUserID uuid.UUID) (Release, error) {
	now := time.Now()
	r := Release{
		ID:           uuid.New(),
		ProjectID:    projectID,
		ReleaseTitle: input.ReleaseTitle,
		ReleaseNotes: input.ReleaseNotes,
		AuthorUserID: authorUserID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := r.Validate(); err != nil {
		return Release{}, err
	}

	return r, nil
}

func (r *Release) Validate() error {
	if r.ReleaseTitle == "" {
		return errReleaseTitleRequired
	}

	return nil
}
