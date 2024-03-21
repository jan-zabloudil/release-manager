package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var AnonUser = &User{}

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	AvatarUrl string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	ReadForToken(ctx context.Context, token string) (User, error)
}

func (s *User) IsAdmin() bool {
	return s.Role.Role() == adminUserRole
}

func (s *User) IsAnon() bool {
	return s == AnonUser
}
