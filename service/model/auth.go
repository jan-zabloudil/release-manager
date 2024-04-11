package model

import (
	"context"

	"github.com/google/uuid"
)

type AuthRepository interface {
	ReadUserIDForToken(ctx context.Context, token string) (uuid.UUID, error)
}
