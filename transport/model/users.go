package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	svcmodel "github.com/jan-zabloudil/release-manager/service/model"
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

func ToAuthUser(u svcmodel.User) AuthUser {
	return AuthUser{
		ID:      u.ID,
		IsAdmin: u.IsAdmin,
	}
}
