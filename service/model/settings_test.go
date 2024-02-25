package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSlackConfigured(t *testing.T) {
	tests := []struct {
		name       string
		slackToken string
		want       bool
	}{
		{"Slack Configured", "some-token", true},
		{"Slack Not Configured", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Settings{
				SlackToken: tt.slackToken,
			}
			assert.Equal(t, tt.want, s.IsSlackConfigured(), "IsSlackConfigured() did not return expected value")
		})
	}
}

func TestIsGithubConfigured(t *testing.T) {
	tests := []struct {
		name        string
		githubToken string
		want        bool
	}{
		{"Github Configured", "some-token", true},
		{"Github Not Configured", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Settings{
				GithubToken: tt.githubToken,
			}
			assert.Equal(t, tt.want, s.IsGithubConfigured(), "IsGithubConfigured() did not return expected value")
		})
	}
}
