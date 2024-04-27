package model

import (
	"testing"

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
			name: "Invalid Project",
			creation: CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Project - invalid github repository url - not absolute url",
			creation: CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
				GithubRepositoryRawURL:    "invalid/url",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project - invalid github repository url - not github host",
			creation: CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
				GithubRepositoryRawURL:    "https://google.com/url/url",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project - invalid github repository url - not repository url #1",
			creation: CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
				GithubRepositoryRawURL:    "https://github.com/owner",
			},
			wantErr: true,
		},
		{
			name: "Invalid Project - invalid github repository url - not repository url #2",
			creation: CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "channelID",
				ReleaseNotificationConfig: ReleaseNotificationConfig{Message: "Test Message"},
				GithubRepositoryRawURL:    "https://github.com/owner/repo/pr",
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

	githubInvalidURL := "invalid/url"
	githubInvalidHost := "https://google.com/url/url"
	githubInvalidRepoURL1 := "https://github.com/owner"
	githubInvalidRepoURL2 := "https://github.com/owner/repo/pr"

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
				Name:           oldName,
				SlackChannelID: oldSlackChannelID,
			},
			update: UpdateProjectInput{
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
			update: UpdateProjectInput{
				Name:           &newInvalidName,
				SlackChannelID: &newSlackChannelID,
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - invalid github repository url - not absolute url",
			project: Project{
				ID:             uuid.New(),
				Name:           oldName,
				SlackChannelID: oldSlackChannelID,
			},
			update: UpdateProjectInput{
				Name:                   &newInvalidName,
				SlackChannelID:         &newSlackChannelID,
				GithubRepositoryRawURL: &githubInvalidURL,
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - invalid github repository url - not github host",
			project: Project{
				ID:             uuid.New(),
				Name:           oldName,
				SlackChannelID: oldSlackChannelID,
			},
			update: UpdateProjectInput{
				Name:                   &newInvalidName,
				SlackChannelID:         &newSlackChannelID,
				GithubRepositoryRawURL: &githubInvalidHost,
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - invalid github repository url - not repository url #1",
			project: Project{
				ID:             uuid.New(),
				Name:           oldName,
				SlackChannelID: oldSlackChannelID,
			},
			update: UpdateProjectInput{
				Name:                   &newInvalidName,
				SlackChannelID:         &newSlackChannelID,
				GithubRepositoryRawURL: &githubInvalidRepoURL1,
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - invalid github repository url - not repository url #2",
			project: Project{
				ID:             uuid.New(),
				Name:           oldName,
				SlackChannelID: oldSlackChannelID,
			},
			update: UpdateProjectInput{
				Name:                   &newInvalidName,
				SlackChannelID:         &newSlackChannelID,
				GithubRepositoryRawURL: &githubInvalidRepoURL2,
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
