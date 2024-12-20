package model

import (
	"errors"
)

var (
	errOrganizationNameRequired  = errors.New("organization name is required")
	errDefaultReleaseMsgRequired = errors.New("default release message is required")
	errSlackMissingToken         = errors.New("token is required to enable slack integration")
	errGithubMissingToken        = errors.New("token is required to enable github integration")
)

type Settings struct {
	OrganizationName      string
	DefaultReleaseMessage string
	Slack                 SlackSettings
	Github                GithubSettings
}

type UpdateSettingsInput struct {
	OrganizationName  *string
	DefaultReleaseMsg *string
	Slack             UpdateSlackSettingsInput
	Github            UpdateGithubSettingsInput
}

type SlackToken string

func (t SlackToken) String() string {
	return string(t)
}

type SlackSettings struct {
	Enabled bool
	Token   SlackToken
}

type UpdateSlackSettingsInput struct {
	Enabled *bool
	Token   *SlackToken
}

type GithubToken string
type GithubWebhookSecret string

func (t GithubToken) String() string {
	return string(t)
}

type GithubSettings struct {
	Enabled       bool
	Token         GithubToken
	WebhookSecret GithubWebhookSecret
}

type UpdateGithubSettingsInput struct {
	Enabled       *bool
	Token         *GithubToken
	WebhookSecret *GithubWebhookSecret
}

func (s *Settings) Update(u UpdateSettingsInput) error {
	if u.OrganizationName != nil {
		s.OrganizationName = *u.OrganizationName
	}
	if u.DefaultReleaseMsg != nil {
		s.DefaultReleaseMessage = *u.DefaultReleaseMsg
	}

	if err := s.Slack.Update(u.Slack); err != nil {
		return err
	}
	if err := s.Github.Update(u.Github); err != nil {
		return err
	}

	return s.Validate()
}

func (s *Settings) Validate() error {
	if s.OrganizationName == "" {
		return errOrganizationNameRequired
	}

	if s.DefaultReleaseMessage == "" {
		return errDefaultReleaseMsgRequired
	}

	if err := s.Slack.Validate(); err != nil {
		return err
	}

	if err := s.Github.Validate(); err != nil {
		return err
	}

	return nil
}

func (s *SlackSettings) Update(u UpdateSlackSettingsInput) error {
	if u.Enabled != nil {
		s.Enabled = *u.Enabled
	}
	if u.Token != nil {
		s.Token = *u.Token
	}

	return s.Validate()
}

func (s *SlackSettings) Validate() error {
	if s.Enabled && s.Token == "" {
		return errSlackMissingToken
	}

	return nil
}

func (s *GithubSettings) Update(u UpdateGithubSettingsInput) error {
	if u.Enabled != nil {
		s.Enabled = *u.Enabled
	}
	if u.Token != nil {
		s.Token = *u.Token
	}
	if u.WebhookSecret != nil {
		s.WebhookSecret = *u.WebhookSecret
	}

	return s.Validate()
}

func (s *GithubSettings) Validate() error {
	if s.Enabled && s.Token == "" {
		return errGithubMissingToken
	}

	return nil
}
