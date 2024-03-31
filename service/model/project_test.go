package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProject_NewProject(t *testing.T) {
	tests := []struct {
		name     string
		creation ProjectCreation
		wantErr  bool
	}{
		{
			name: "Valid Project",
			creation: ProjectCreation{
				Name:                      "Test Project",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
			},
			wantErr: false,
		},
		{
			name: "Invalid Project",
			creation: ProjectCreation{
				Name:                      "",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProject(tt.creation)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProject_Update(t *testing.T) {
	oldName := "Old Name"
	oldSlackChannelID := "oldChannelID"
	newValidName := "New Name"
	newInvalidName := ""
	newSlackChannelID := "newChannelID"

	tests := []struct {
		name    string
		project Project
		update  ProjectUpdate
		wantErr bool
	}{
		{
			name: "Valid Update",
			project: Project{
				ID:             uuid.New(),
				Name:           oldName,
				SlackChannelID: oldSlackChannelID,
			},
			update: ProjectUpdate{
				Name:           &newValidName,
				SlackChannelID: &newSlackChannelID,
			},
			wantErr: false,
		},
		{
			name: "Invalid Update",
			project: Project{
				ID:             uuid.New(),
				Name:           oldName,
				SlackChannelID: oldSlackChannelID,
			},
			update: ProjectUpdate{
				Name:           &newInvalidName,
				SlackChannelID: &newSlackChannelID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Update(tt.update)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		wantErr bool
	}{
		{
			name: "Valid Project",
			project: Project{
				ID:             uuid.New(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
			},
			wantErr: false,
		},
		{
			name: "Invalid Project",
			project: Project{
				ID:             uuid.New(),
				Name:           "",
				SlackChannelID: "channelID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
