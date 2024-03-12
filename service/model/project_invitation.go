package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProjectInvitationRepository interface {
	Insert(ctx context.Context, projectID uuid.UUID, email string, role ProjectRole, invitedByUserID uuid.UUID) (ProjectInvitation, error)
	ReadByEmail(ctx context.Context, projectID uuid.UUID, email string) (ProjectInvitation, error)
	Read(ctx context.Context, invitationID uuid.UUID) (ProjectInvitation, error)
	ReadAll(ctx context.Context, projectID uuid.UUID) ([]ProjectInvitation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectInvitation struct {
	ID              uuid.UUID
	ProjectID       uuid.UUID
	Email           string
	Role            ProjectRole
	InvitedByUserID uuid.UUID
	CreatedAt       time.Time
}
