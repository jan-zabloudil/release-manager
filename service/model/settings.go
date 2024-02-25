package model

import (
	"context"
)

type SettingsRepository interface {
	Set(ctx context.Context, s Settings) (Settings, error)
	Read(ctx context.Context) (Settings, error)
}

type Settings struct {
	OrganizationName  string
	SlackToken        string
	GithubToken       string
	DefaultReleaseMsg string
}

func (s Settings) IsSlackConfigured() bool {
	return s.SlackToken != ""
}

func (s Settings) IsGithubConfigured() bool {
	return s.GithubToken != ""
}
