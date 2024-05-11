package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	AvatarURL string    `db:"avatar_url"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToSvcUser(u User) svcmodel.User {
	return svcmodel.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Role:      svcmodel.UserRole(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToSvcUsers(users []User) []svcmodel.User {
	u := make([]svcmodel.User, 0, len(users))

	for _, user := range users {
		u = append(u, ToSvcUser(user))
	}

	return u
}
