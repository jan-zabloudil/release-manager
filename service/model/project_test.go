package model

import (
	"testing"

	"release-manager/pkg/pointer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProject_NewProject(t *testing.T) {
	tests := []struct {
		name     string
		creation CreateProjectInput
		wantErr  bool
	}{
		{
			name: "Valid Project",
			creation: CreateProjectInput{
				Name:                      "Test Project",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			creation: CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
			},
			wantErr: true,
		},
		{
			name: "Missing release notification config message",
			creation: CreateProjectInput{
				Name:                      "Test project",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: ""},
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
	tests := []struct {
		name    string
		project Project
		update  UpdateProjectInput
		wantErr bool
	}{
		{
			name: "Valid Update",
			project: Project{
				ID:             uuid.New(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "Test Message",
				},
			},
			update: UpdateProjectInput{
				Name:           pointer.StringPtr("New name"),
				SlackChannelID: pointer.StringPtr("newChannelID"),
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			project: Project{
				ID:             uuid.New(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "Test Message",
				},
			},
			update: UpdateProjectInput{
				Name:           pointer.StringPtr(""),
				SlackChannelID: pointer.StringPtr("newChannelID"),
			},
			wantErr: true,
		},
		{
			name: "Missing release config message",
			project: Project{
				ID:             uuid.New(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "Test Message",
				},
			},
			update: UpdateProjectInput{
				Name:           pointer.StringPtr("New name"),
				SlackChannelID: pointer.StringPtr("newChannelID"),
				ReleaseNotificationConfigUpdate: UpdateReleaseNotificationConfigInput{
					Message: pointer.StringPtr(""),
				},
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
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "Test Message",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			project: Project{
				ID:             uuid.New(),
				Name:           "",
				SlackChannelID: "channelID",
			},
			wantErr: true,
		},
		{
			name: "Missing release notification config message",
			project: Project{
				ID:             uuid.New(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "",
				},
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

func TestReleaseNotificationConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ReleaseNotificationConfig
		wantErr bool
	}{
		{
			name: "Valid Config",
			config: ReleaseNotificationConfig{
				Message: "Test Message",
			},
			wantErr: false,
		},
		{
			name: "Invalid Config",
			config: ReleaseNotificationConfig{
				Message: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReleaseNotificationConfig_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		config *ReleaseNotificationConfig
		want   bool
	}{
		{
			name:   "Empty Config",
			config: &ReleaseNotificationConfig{},
			want:   true,
		},
		{
			name: "Non-Empty Config",
			config: &ReleaseNotificationConfig{
				Message: "Test Message",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsEmpty()
			assert.Equal(t, tt.want, got)
		})
	}
}
