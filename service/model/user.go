package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	AvatarURL string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToUser(id uuid.UUID, email, name, avatarURL, roleStr string, createdAt, updatedAt time.Time) (User, error) {
	role, err := NewUserRole(roleStr)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        id,
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s User) HasAtLeastRole(role UserRole) bool {
	return s.Role.IsRoleAtLeast(role)
}
