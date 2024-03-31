package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestUserRole_NewUserRole(t *testing.T) {
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
			_, err := NewUserRole(tc.roleStr)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
