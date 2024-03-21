package model

import (
	"testing"

	svcerr "release-manager/service/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSCMRepo_NewSCMRepo(t *testing.T) {
	appID := uuid.New()
	testCases := []struct {
		name      string
		platform  string
		repoURL   string
		wantError error
	}{
		{
			name:      "Valid GitHub Repo",
			platform:  "github",
			repoURL:   "https://github.com/owner/repo",
			wantError: nil,
		},
		{
			name:      "Invalid Platform",
			platform:  "invalid",
			repoURL:   "https://github.com/owner/repo",
			wantError: svcerr.ErrUnknownSCMRepoPlatform,
		},
		{
			name:      "Invalid GitHub Repo URL",
			platform:  "github",
			repoURL:   "https://github.com/owner",
			wantError: svcerr.ErrInvalidGithubRepoUrl,
		},
		{
			name:      "Non-Absolute GitHub Repo URL",
			platform:  "github",
			repoURL:   "/owner/repo",
			wantError: svcerr.ErrInvalidSCMRepoURL,
		},
		{
			name:      "Non-GitHub Host URL",
			platform:  "github",
			repoURL:   "https://gitlab.com/owner/repo",
			wantError: svcerr.ErrInvalidGithubHostUrl,
		},
		{
			name:      "Empty Platform",
			platform:  "",
			repoURL:   "https://github.com/owner/repo",
			wantError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSCMRepo(appID, tc.platform, tc.repoURL)
			if tc.wantError != nil {
				assert.ErrorIs(t, err, tc.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSCMRepo_NewEmptySCMRepo(t *testing.T) {
	repo := NewEmptySCMRepo()
	assert.False(t, repo.IsSet())
}

func TestSCMRepo_IsSet(t *testing.T) {
	appID := uuid.New()
	repoURL := "https://github.com/owner/repo"

	testCases := []struct {
		name     string
		repo     SCMRepo
		expected bool
	}{
		{
			name:     "githubRepo IsSet",
			repo:     &githubRepo{appID: appID, repoURL: repoURL},
			expected: true,
		},
		{
			name:     "emptyRepo IsSet",
			repo:     &emptyRepo{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.repo.IsSet())
		})
	}
}
