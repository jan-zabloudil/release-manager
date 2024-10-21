package model

import (
	"testing"

	"release-manager/pkg/id"

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
				ID:        id.User{},
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

func TestUser_IsAdmin(t *testing.T) {
	testCases := []struct {
		name     string
		userRole UserRole
		expected bool
	}{
		{"Admin user", UserRoleAdmin, true},
		{"Non-admin user", UserRoleUser, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := User{
				ID:   id.User{},
				Role: tc.userRole,
			}

			assert.Equal(t, tc.expected, user.IsAdmin())
		})
	}
}
