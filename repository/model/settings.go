package model

import (
	svcmodel "release-manager/service/model"
)

const (
	orgNameKey           = "organization_name"
	slackTokenKey        = "slack_token"
	githubTokenKey       = "github_access_token"
	defaultReleaseMsgKey = "default_release_msg"
)

type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ToDBSettings(orgName, slackTkn, githubTkn, rlsMsg string) []Setting {
	return []Setting{
		{Key: orgNameKey, Value: orgName},
		{Key: slackTokenKey, Value: slackTkn},
		{Key: githubTokenKey, Value: githubTkn},
		{Key: defaultReleaseMsgKey, Value: rlsMsg},
	}
}

func ToSvcSettings(settings []Setting) (s svcmodel.Settings) {
	for _, setting := range settings {
		switch setting.Key {
		case slackTokenKey:
			s.SlackToken = setting.Value
		case orgNameKey:
			s.OrganizationName = setting.Value
		case githubTokenKey:
			s.GithubToken = setting.Value
		case defaultReleaseMsgKey:
			s.DefaultReleaseMsg = setting.Value
		}
	}

	return s
}
