package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsAdmin(t *testing.T) {
	adminRole, _ := NewUserRole(adminUserRole)
	basicRole := NewBasicUserRole()

	tests := []struct {
		name string
		user User
		want bool
	}{
		{
			name: "admin user",
			user: User{
				Role: adminRole,
			},
			want: true,
		},
		{
			name: "basic user",
			user: User{
				Role: basicRole,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.user.IsAdmin(), "User.IsAdmin() did not return the expected value")
		})
	}
}

func TestIsAnon(t *testing.T) {
	user := User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		AvatarURL: "https://example.com/avatar.jpg",
		Role:      NewBasicUserRole(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name string
		user *User
		want bool
	}{
		{
			name: "anonymous user",
			user: AnonUser,
			want: true,
		},
		{
			name: "authenticated user",
			user: &user,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.user.IsAnon(), "User.IsAnon() did not return the expected value")
		})
	}
}
