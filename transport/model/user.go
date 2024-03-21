package model

import (
	"context"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type UserService interface {
	GetForToken(ctx context.Context, token string) (svcmodel.User, error)
	Get(ctx context.Context, id uuid.UUID) (svcmodel.User, error)
	GetAll(ctx context.Context) ([]svcmodel.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarUrl string    `json:"avatar_url"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToNetUser(id uuid.UUID, role, email, name, avatarUrl string, createdAt, updatedAt time.Time) User {
	return User{
		ID:        id,
		Role:      role,
		Email:     email,
		Name:      name,
		AvatarUrl: avatarUrl,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func ToNetUsers(users []svcmodel.User) []User {
	u := make([]User, 0, len(users))
	for _, user := range users {
		u = append(u, ToNetUser(
			user.ID,
			user.Role.Role(),
			user.Email,
			user.Name,
			user.AvatarUrl,
			user.CreatedAt,
			user.UpdatedAt,
		))
	}

	return u
}
