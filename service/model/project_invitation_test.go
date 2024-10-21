package model

import (
	"testing"

	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/id"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectInvitation_NewProjectInvitation(t *testing.T) {
	tests := []struct {
		name     string
		creation CreateProjectInvitationInput
		wantErr  bool
	}{
		{
			name: "Valid Project Invitation",
			creation: CreateProjectInvitationInput{
				ProjectID:   uuid.New(),
				Email:       "test@example.com",
				ProjectRole: "editor",
			},
			wantErr: false,
		},
		{
			name: "Invalid Project Invitation - missing email",
			creation: CreateProjectInvitationInput{
				ProjectID:   uuid.New(),
				Email:       "test@example.com",
				ProjectRole: "owner",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project Invitation - missing email",
			creation: CreateProjectInvitationInput{
				ProjectID:   uuid.New(),
				Email:       "",
				ProjectRole: "viewer",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project Invitation - invalid email",
			creation: CreateProjectInvitationInput{
				ProjectID:   uuid.New(),
				Email:       "test@test",
				ProjectRole: "viewer",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project Invitation - invalid role",
			creation: CreateProjectInvitationInput{
				ProjectID:   uuid.New(),
				Email:       "test@test.tt",
				ProjectRole: "admin",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn, err := cryptox.NewToken()
			if err != nil {
				t.Fatal(err)
			}

			_, err = NewProjectInvitation(tt.creation, tkn, id.AuthUser{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectInvitation_Validate(t *testing.T) {
	tests := []struct {
		name       string
		invitation ProjectInvitation
		wantErr    bool
	}{
		{
			name: "Valid Project Invitation",
			invitation: ProjectInvitation{
				ID:          id.NewProjectInvitation(),
				ProjectID:   uuid.New(),
				Email:       "test@example.com",
				ProjectRole: "editor",
				Status:      "pending",
			},
			wantErr: false,
		},
		{
			name: "Invalid Project Invitation - missing email",
			invitation: ProjectInvitation{
				ID:          id.NewProjectInvitation(),
				ProjectID:   uuid.New(),
				Email:       "",
				ProjectRole: "editor",
				Status:      "pending",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project Invitation - invalid email",
			invitation: ProjectInvitation{
				ID:          id.NewProjectInvitation(),
				ProjectID:   uuid.New(),
				Email:       "test@test",
				ProjectRole: "editor",
				Status:      "pending",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project Invitation - invalid role",
			invitation: ProjectInvitation{
				ID:          id.NewProjectInvitation(),
				ProjectID:   uuid.New(),
				Email:       "test@test",
				ProjectRole: "admin",
				Status:      "pending",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project Invitation - invalid status",
			invitation: ProjectInvitation{
				ID:          id.NewProjectInvitation(),
				ProjectID:   uuid.New(),
				Email:       "test@test",
				ProjectRole: "owner",
				Status:      "new",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.invitation.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectInvitation_Accept(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus ProjectInvitationStatus
		expectedError error
	}{
		{
			name:          "invitation is already accepted",
			initialStatus: InvitationStatusAccepted,
			expectedError: ErrProjectInvitationAlreadyAccepted,
		},
		{
			name:          "invitation is pending",
			initialStatus: InvitationStatusPending,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitation := &ProjectInvitation{
				Status: tt.initialStatus,
			}

			err := invitation.Accept()

			assert.Equal(t, tt.expectedError, err)

			if tt.expectedError == nil {
				assert.Equal(t, InvitationStatusAccepted, invitation.Status)
				assert.False(t, invitation.UpdatedAt.IsZero(), "UpdatedAt should be set")
			}
		})
	}
}
