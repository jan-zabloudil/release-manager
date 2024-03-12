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

var AnonUser = &User{}

func (u *User) IsAnon() bool {
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

func ToNetUser(
	id uuid.UUID,
	isAdmin bool,
	email, name, avatarUrl string,
	createdAt, updatedAt time.Time,
) User {
	return User{
		ID:        id,
		IsAdmin:   isAdmin,
		Email:     email,
		Name:      name,
		AvatarUrl: avatarUrl,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func ToSvcUser(u User) svcmodel.User {
	return svcmodel.User(u)
}
