package model

import (
	"context"

	"github.com/google/uuid"
)

const (
	InvitationSentStatus = "invited"
	MemberCreatedStatus  = "created"
)

type UserService interface {
	GetByEmail(ctx context.Context, email string) (User, error)
}

type ProjectInvitationService interface {
	GetByEmail(ctx context.Context, projectID uuid.UUID, email string) (ProjectInvitation, error)
	Create(ctx context.Context, projectID uuid.UUID, email string, role ProjectRole, invitedByUserID uuid.UUID) (ProjectInvitation, error)
}

type ProjectMemberService interface {
	Get(ctx context.Context, projectID, userID uuid.UUID) (ProjectMember, error)
	Create(ctx context.Context, projectID, userID uuid.UUID, role ProjectRole, invitedByUserID uuid.UUID) (ProjectMember, error)
}

type ProjectMembershipRequest struct {
	ProjectID         uuid.UUID
	Email             string
	Role              ProjectRole
	RequestedByUserID uuid.UUID
}

type ProjectMembershipResponse struct {
	Status   string
	Resource any
}
