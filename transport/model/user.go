package model

import (
	"context"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type UserService interface {
	GetForToken(ctx context.Context, token string) (svcmodel.User, error)
}

var AnonUser = &AuthUser{}

type AuthUser struct {
	ID      uuid.UUID
	IsAdmin bool
}

func (u *AuthUser) IsAnon() bool {
	return u == AnonUser
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarUrl string    `json:"avatar_url"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAuthUser(id uuid.UUID, isAdmin bool) AuthUser {
	return AuthUser{
		ID:      id,
		IsAdmin: isAdmin,
	}
}
