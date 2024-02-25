package model

import (
	"context"

	svcmodel "release-manager/service/model"
)

type SettingsService interface {
	Set(ctx context.Context, s svcmodel.Settings) (svcmodel.Settings, error)
	Get(ctx context.Context) (svcmodel.Settings, error)
}

type Settings struct {
	OrganizationName  *string `json:"organization_name"`
	SlackToken        *string `json:"slack_token"`
	GithubToken       *string `json:"github_token"`
	DefaultReleaseMsg *string `json:"default_release_msg"`
}

func ToSvcSettings(s svcmodel.Settings, orgName, slackTkn, githubTkn, rlsMsg *string) svcmodel.Settings {
	if orgName != nil {
		s.OrganizationName = *orgName
	}
	if slackTkn != nil {
		s.SlackToken = *slackTkn
	}
	if githubTkn != nil {
		s.GithubToken = *githubTkn
	}
	if rlsMsg != nil {
		s.DefaultReleaseMsg = *rlsMsg
	}

	return s
}

func ToNetSettings(orgName, slackTkn, githubTkn, rlsMsg string) Settings {
	return Settings{
		OrganizationName:  &orgName,
		SlackToken:        &slackTkn,
		GithubToken:       &githubTkn,
		DefaultReleaseMsg: &rlsMsg,
	}
}
