package model

import (
	"testing"

	"release-manager/pkg/id"
	"release-manager/pkg/pointer"

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
	slackChannelID := SlackChannelID("newChannelID")

	tests := []struct {
		name    string
		project Project
		update  UpdateProjectInput
		wantErr bool
	}{
		{
			name: "Valid Update",
			project: Project{
				ID:             id.NewProject(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "Test Message",
				},
			},
			update: UpdateProjectInput{
				Name:           pointer.StringPtr("New name"),
				SlackChannelID: &slackChannelID,
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			project: Project{
				ID:             id.NewProject(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "Test Message",
				},
			},
			update: UpdateProjectInput{
				Name:           pointer.StringPtr(""),
				SlackChannelID: &slackChannelID,
			},
			wantErr: true,
		},
		{
			name: "Missing release config message",
			project: Project{
				ID:             id.NewProject(),
				Name:           "Test Project",
				SlackChannelID: "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{
					Message: "Test Message",
				},
			},
			update: UpdateProjectInput{
				Name:           pointer.StringPtr("New name"),
				SlackChannelID: &slackChannelID,
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
				ID:             id.NewProject(),
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
				ID:             id.NewProject(),
				Name:           "",
				SlackChannelID: "channelID",
			},
			wantErr: true,
		},
		{
			name: "Missing release notification config message",
			project: Project{
				ID:             id.NewProject(),
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

func TestProject_GithubOwnerSlug(t *testing.T) {
	tests := []struct {
		name           string
		project        *Project
		expectedResult *string
	}{
		{
			name:           "GithubRepo is nil",
			project:        &Project{GithubRepo: nil},
			expectedResult: nil,
		},
		{
			name:           "OwnerSlug is returned",
			project:        &Project{GithubRepo: &GithubRepo{OwnerSlug: "owner123"}},
			expectedResult: pointer.StringPtr("owner123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.project.GithubOwnerSlug()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestProject_GithubRepoSlug(t *testing.T) {
	tests := []struct {
		name           string
		project        *Project
		expectedResult *string
	}{
		{
			name:           "GithubRepo is nil",
			project:        &Project{GithubRepo: nil},
			expectedResult: nil,
		},
		{
			name:           "RepoSlug is returned",
			project:        &Project{GithubRepo: &GithubRepo{RepoSlug: "repo123"}},
			expectedResult: pointer.StringPtr("repo123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.project.GithubRepoSlug()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
