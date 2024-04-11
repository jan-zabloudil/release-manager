package model

import (
	"testing"
	"time"

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

func TestUser_ToUser(t *testing.T) {
	testCases := []struct {
		name    string
		roleStr string
		wantErr bool
	}{
		{"Valid role - admin", "admin", false},
		{"Valid role - user", "user", false},
		{"Invalid role", "invalid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ToUser(uuid.New(), "test@example.com", "Test User", "https://example.com/avatar.jpg", tc.roleStr, time.Now(), time.Now())
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
