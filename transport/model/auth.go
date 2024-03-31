package model

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
