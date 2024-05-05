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

func TestProjectMember_UpdateProjectRole(t *testing.T) {
	tests := []struct {
		name    string
		role    ProjectRole
		newRole ProjectRole
		wantErr bool
	}{
		{
			name:    "Update Role from Editor to Viewer",
			role:    ProjectRoleEditor,
			newRole: ProjectRoleViewer,
			wantErr: false,
		},
		{
			name:    "Update Role from Viewer to Owner",
			role:    ProjectRoleViewer,
			newRole: ProjectRoleOwner,
			wantErr: true,
		},
		{
			name:    "Update Role to invalid role",
			role:    ProjectRoleViewer,
			newRole: "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member := ProjectMember{
				User:        User{ID: uuid.New(), Email: "test@example.com"},
				ProjectID:   uuid.New(),
				ProjectRole: tt.role,
			}
			err := member.UpdateProjectRole(tt.newRole)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newRole, member.ProjectRole)
			}
		})
	}
}

func TestProjectMember_HasAtLeastProjectRole(t *testing.T) {
	tests := []struct {
		name            string
		role            ProjectRole
		wantAtLeastRole ProjectRole
		want            bool
	}{
		{
			name:            "Editor, want at least Viewer",
			role:            ProjectRoleEditor,
			wantAtLeastRole: ProjectRoleViewer,
			want:            true,
		},
		{
			name:            "Viewer, want at least Viewer",
			role:            ProjectRoleViewer,
			wantAtLeastRole: ProjectRoleViewer,
			want:            true,
		},
		{
			name:            "Viewer, want at least Editor",
			role:            ProjectRoleViewer,
			wantAtLeastRole: ProjectRoleEditor,
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member := ProjectMember{
				ProjectRole: tt.role,
			}
			got := member.HasAtLeastProjectRole(tt.wantAtLeastRole)
			assert.Equal(t, tt.want, got)
		})
	}
}
