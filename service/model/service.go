package model

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	AuthorizeAdminRole(ctx context.Context, userID uuid.UUID) error
	AuthorizeRole(ctx context.Context, userID uuid.UUID, role UserRole) error
}
