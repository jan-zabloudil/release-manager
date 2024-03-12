package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProjectMemberRepository interface {
	Insert(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, role ProjectRole, invitedByUserID uuid.UUID) (ProjectMember, error)
	Read(ctx context.Context, projectID, userID uuid.UUID) (ProjectMember, error)
	Delete(ctx context.Context, projectID, userID uuid.UUID) error
	ReadAll(ctx context.Context, projectID uuid.UUID) ([]ProjectMember, error)
	Update(ctx context.Context, member ProjectMember) (ProjectMember, error)
}

type ProjectMember struct {
	User            User
	ProjectID       uuid.UUID
	Role            ProjectRole
	InvitedByUserId uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
