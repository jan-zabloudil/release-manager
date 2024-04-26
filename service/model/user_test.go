package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_HasAtLeastRole(t *testing.T) {
	testCases := []struct {
		name     string
		userRole UserRole
		testRole UserRole
		expected bool
	}{
		{"Admin has at least User role", UserRoleAdmin, UserRoleUser, true},
		{"Admin has at least Admin role", UserRoleAdmin, UserRoleAdmin, true},
		{"User does not have at least Admin role", UserRoleUser, UserRoleAdmin, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := User{
				ID:        uuid.New(),
				Email:     "test@example.com",
				Name:      "Test User",
				AvatarURL: "https://example.com/avatar.jpg",
				Role:      tc.userRole,
			}

			assert.Equal(t, tc.expected, user.HasAtLeastRole(tc.testRole))
		})
	}
}

func TestUserRole_IsRoleAtLeast(t *testing.T) {
	testCases := []struct {
		name     string
		userRole UserRole
		testRole UserRole
		expected bool
	}{
		{"Admin has at least User role", UserRoleAdmin, UserRoleUser, true},
		{"Admin has at least Admin role", UserRoleAdmin, UserRoleAdmin, true},
		{"User does not have at least Admin role", UserRoleUser, UserRoleAdmin, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.userRole.IsRoleAtLeast(tc.testRole))
		})
	}
}
