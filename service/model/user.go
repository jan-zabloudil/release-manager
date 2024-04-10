package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	ReadForToken(ctx context.Context, token string) (User, error)
	Read(ctx context.Context, id uuid.UUID) (User, error)
	ReadAll(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

var AnonUser = &User{}

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	AvatarURL string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *User) IsAdmin() bool {
	return s.Role.Role() == adminUserRole
}

func (s *User) IsAnon() bool {
	return s == AnonUser
}
