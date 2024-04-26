package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUser(u svcmodel.User) User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt.Local(),
		UpdatedAt: u.UpdatedAt.Local(),
	}
}

func ToUsers(users []svcmodel.User) []User {
	u := make([]User, 0, len(users))
	for _, user := range users {
		u = append(u, ToUser(user))
	}

	return u
}
