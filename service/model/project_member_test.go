package model

import (
	"testing"

	"release-manager/pkg/id"

	"github.com/stretchr/testify/assert"
)

func TestProjectMember_NewProjectMember(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		role    ProjectRole
		wantErr bool
	}{
		{
			name:    "Valid Project Member - Owner",
			user:    User{ID: id.User{}, Email: "test@example.com"},
			role:    ProjectRoleOwner,
			wantErr: false,
		},
		{
			name:    "Valid Project Member - Editor",
			user:    User{ID: id.User{}, Email: "test@example.com"},
			role:    ProjectRoleEditor,
			wantErr: false,
		},
		{
			name:    "Valid Project Member - Viewer",
			user:    User{ID: id.User{}, Email: "test@example.com"},
			role:    ProjectRoleViewer,
			wantErr: false,
		},
		{
			name:    "Invalid Project Member - Invalid Role",
			user:    User{ID: id.User{}, Email: "test@example.com"},
			role:    "admin",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProjectMember(tt.user, id.NewProject(), tt.role)
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
				User:        User{ID: id.User{}, Email: "test@example.com"},
				ProjectID:   id.NewProject(),
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
		member          ProjectMember
		wantAtLeastRole ProjectRole
		want            bool
	}{
		{
			name: "Editor, want at least Viewer",
			member: ProjectMember{
				ProjectRole: ProjectRoleEditor,
				User: User{
					Role: UserRoleUser,
				},
			},
			wantAtLeastRole: ProjectRoleViewer,
			want:            true,
		},
		{
			name: "Viewer, want at least Viewer",
			member: ProjectMember{
				ProjectRole: ProjectRoleViewer,
				User: User{
					Role: UserRoleUser,
				},
			},
			wantAtLeastRole: ProjectRoleViewer,
			want:            true,
		},
		{
			name: "Viewer, want at least Editor",
			member: ProjectMember{
				ProjectRole: ProjectRoleViewer,
				User: User{
					Role: UserRoleUser,
				},
			},
			wantAtLeastRole: ProjectRoleEditor,
			want:            false,
		},
		{
			name: "Viewer (but also admin user), want at least Editor",
			member: ProjectMember{
				ProjectRole: ProjectRoleViewer,
				User: User{
					Role: UserRoleAdmin,
				},
			},
			wantAtLeastRole: ProjectRoleEditor,
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.member.SatisfiesRequiredRole(tt.wantAtLeastRole)
			assert.Equal(t, tt.want, got)
		})
	}
}
