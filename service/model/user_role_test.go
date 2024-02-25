package model

import (
	"testing"

	svcerr "release-manager/service/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserRole(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		expectedErr error
	}{
		{
			name:        "valid user role",
			role:        basicUserRole,
			expectedErr: nil,
		},
		{
			name:        "valid admin role",
			role:        adminUserRole,
			expectedErr: nil,
		},
		{
			name:        "invalid role",
			role:        "invalidRole",
			expectedErr: svcerr.ErrInvalidUserRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewUserRole(tt.role)

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.role, result.Role())
			}
		})
	}
}

func TestNewBasicUserRole(t *testing.T) {
	result := NewBasicUserRole()
	assert.Equal(t, basicUserRole, result.Role())
}
