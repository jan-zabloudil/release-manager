package model

import (
	svcmodel "release-manager/service/model"
)

type UpdateSettingsInput struct {
	OrganizationName      *string                   `json:"organization_name" validate:"omitempty,min=1"`
	DefaultReleaseMessage *string                   `json:"default_release_message" validate:"omitempty,min=1"`
	Slack                 UpdateSlackSettingsInput  `json:"slack"`
	Github                UpdateGithubSettingsInput `json:"github"`
}

type UpdateSlackSettingsInput struct {
	Enabled *bool                `json:"enabled"`
	Token   *svcmodel.SlackToken `json:"token"`
}

type UpdateGithubSettingsInput struct {
	Enabled       *bool                         `json:"enabled"`
	Token         *svcmodel.GithubToken         `json:"token"`
	WebhookSecret *svcmodel.GithubWebhookSecret `json:"webhook_secret"`
}

type Settings struct {
	OrganizationName      string         `json:"organization_name"`
	DefaultReleaseMessage string         `json:"default_release_message"`
	Slack                 SlackSettings  `json:"slack"`
	Github                GithubSettings `json:"github"`
}

type SlackSettings struct {
	Enabled bool                `json:"enabled"`
	Token   svcmodel.SlackToken `json:"token"`
}

type GithubSettings struct {
	Enabled       bool                         `json:"enabled"`
	Token         svcmodel.GithubToken         `json:"token"`
	WebhookSecret svcmodel.GithubWebhookSecret `json:"webhook_secret"`
}

func ToSvcUpdateSettingsInput(u UpdateSettingsInput) svcmodel.UpdateSettingsInput {
	return svcmodel.UpdateSettingsInput{
		OrganizationName:  u.OrganizationName,
		DefaultReleaseMsg: u.DefaultReleaseMessage,
		Slack: svcmodel.UpdateSlackSettingsInput{
			Enabled: u.Slack.Enabled,
			Token:   u.Slack.Token,
		},
		Github: svcmodel.UpdateGithubSettingsInput{
			Enabled:       u.Github.Enabled,
			Token:         u.Github.Token,
			WebhookSecret: u.Github.WebhookSecret,
		},
	}
}

func ToSettings(s svcmodel.Settings) Settings {
	return Settings{
		OrganizationName:      s.OrganizationName,
		DefaultReleaseMessage: s.DefaultReleaseMessage,
		Slack: SlackSettings{
			Enabled: s.Slack.Enabled,
			Token:   s.Slack.Token,
		},
		Github: GithubSettings{
			Enabled:       s.Github.Enabled,
			Token:         s.Github.Token,
			WebhookSecret: s.Github.WebhookSecret,
		},
	}
}
