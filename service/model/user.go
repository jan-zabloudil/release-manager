package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	AvatarUrl string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	ReadForToken(ctx context.Context, token string) (User, error)
	Read(ctx context.Context, id uuid.UUID) (User, error)
	ReadByEmail(ctx context.Context, email string) (User, error)
}
