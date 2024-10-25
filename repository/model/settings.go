package model

import (
	"encoding/json"
	"fmt"

	svcmodel "release-manager/service/model"
)

const (
	keyOrganizationName      = "organization_name"
	keyDefaultReleaseMessage = "default_release_message"
	keySlack                 = "slack"
	keyGithub                = "github"
)

// SettingsValue represents a key-value pair for settings in the database table.
type SettingsValue struct {
	Key   string          `db:"key"`
	Value json.RawMessage `db:"value"`
}

type SlackSettings struct {
	Enabled bool                `json:"enabled"`
	Token   svcmodel.SlackToken `json:"token"`
}

type GithubSettings struct {
	Enabled       bool                 `json:"enabled"`
	Token         svcmodel.GithubToken `json:"token"`
	WebhookSecret string               `json:"webhook_secret"`
}

func ToSettingsValues(s svcmodel.Settings) ([]SettingsValue, error) {
	var sv []SettingsValue

	orgName, err := toSettingsValue(keyOrganizationName, s.OrganizationName)
	if err != nil {
		return nil, err
	}
	rlsMessage, err := toSettingsValue(keyDefaultReleaseMessage, s.DefaultReleaseMessage)
	if err != nil {
		return nil, err
	}
	slack, err := toSettingsValue(keySlack, SlackSettings(s.Slack))
	if err != nil {
		return nil, err
	}
	github, err := toSettingsValue(keyGithub, GithubSettings(s.Github))
	if err != nil {
		return nil, err
	}

	return append(sv, orgName, rlsMessage, slack, github), nil
}

func toSettingsValue(key string, v any) (SettingsValue, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return SettingsValue{}, err
	}

	return SettingsValue{Key: key, Value: data}, nil
}

func ToSvcSettings(sv []SettingsValue) (svcmodel.Settings, error) {
	var s svcmodel.Settings
	var slackSettings SlackSettings
	var githubSettings GithubSettings

	for _, settingsValue := range sv {
		switch settingsValue.Key {
		case keyOrganizationName:
			if err := json.Unmarshal(settingsValue.Value, &s.OrganizationName); err != nil {
				return svcmodel.Settings{}, err
			}
		case keyDefaultReleaseMessage:
			if err := json.Unmarshal(settingsValue.Value, &s.DefaultReleaseMessage); err != nil {
				return svcmodel.Settings{}, err
			}
		case keySlack:
			if err := json.Unmarshal(settingsValue.Value, &slackSettings); err != nil {
				return svcmodel.Settings{}, err
			}
		case keyGithub:
			if err := json.Unmarshal(settingsValue.Value, &githubSettings); err != nil {
				return svcmodel.Settings{}, err
			}
		default:
			return svcmodel.Settings{}, fmt.Errorf("unknown key: %s", settingsValue.Key)
		}
	}

	s.Slack = svcmodel.SlackSettings(slackSettings)
	s.Github = svcmodel.GithubSettings(githubSettings)

	return s, nil
}
