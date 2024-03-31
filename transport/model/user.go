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

func ToUser(id uuid.UUID, role svcmodel.UserRole, email, name, avatarURL string, createdAt, updatedAt time.Time) User {
	return User{
		ID:        id,
		Role:      string(role),
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
		CreatedAt: createdAt.Local(),
		UpdatedAt: updatedAt.Local(),
	}
}

func ToUsers(users []svcmodel.User) []User {
	u := make([]User, 0, len(users))
	for _, user := range users {
		u = append(u, ToUser(
			user.ID,
			user.Role,
			user.Email,
			user.Name,
			user.AvatarURL,
			user.CreatedAt,
			user.UpdatedAt,
		))
	}

	return u
}
