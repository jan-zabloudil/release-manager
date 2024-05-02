package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectMember_NewProjectMember(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		projectID uuid.UUID
		role      ProjectRole
		wantErr   bool
	}{
		{
			name:      "Valid Project Member - Owner",
			user:      User{ID: uuid.New(), Email: "test@example.com"},
			projectID: uuid.New(),
			role:      ProjectRoleOwner,
			wantErr:   false,
		},
		{
			name:      "Valid Project Member - Editor",
			user:      User{ID: uuid.New(), Email: "test@example.com"},
			projectID: uuid.New(),
			role:      ProjectRoleEditor,
			wantErr:   false,
		},
		{
			name:      "Valid Project Member - Viewer",
			user:      User{ID: uuid.New(), Email: "test@example.com"},
			projectID: uuid.New(),
			role:      ProjectRoleViewer,
			wantErr:   false,
		},
		{
			name:      "Invalid Project Member - Invalid Role",
			user:      User{ID: uuid.New(), Email: "test@example.com"},
			projectID: uuid.New(),
			role:      "admin",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProjectMember(tt.user, tt.projectID, tt.role)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
