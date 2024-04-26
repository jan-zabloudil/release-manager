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

type Settings struct {
	OrganizationName      string
	DefaultReleaseMessage string
	Slack                 SlackSettings
	Github                GithubSettings
}

type SlackSettings struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
}

type GithubSettings struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
}

func (s *Settings) MarshalJSON() ([]byte, error) {
	type KeyValue struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}

	return json.Marshal([]KeyValue{
		{Key: keyOrganizationName, Value: s.OrganizationName},
		{Key: keyDefaultReleaseMessage, Value: s.DefaultReleaseMessage},
		{Key: keySlack, Value: s.Slack},
		{Key: keyGithub, Value: s.Github},
	})
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	type KeyRawValue struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}

	var pairs []KeyRawValue
	if err := json.Unmarshal(data, &pairs); err != nil {
		return err
	}

	for _, kv := range pairs {
		switch kv.Key {
		case keyOrganizationName:
			if err := json.Unmarshal(kv.Value, &s.OrganizationName); err != nil {
				return fmt.Errorf("cannot unmarshal %s: %w", kv.Key, err)
			}
		case keyDefaultReleaseMessage:
			if err := json.Unmarshal(kv.Value, &s.DefaultReleaseMessage); err != nil {
				return fmt.Errorf("cannot unmarshal %s: %w", kv.Key, err)
			}
		case keySlack:
			if err := json.Unmarshal(kv.Value, &s.Slack); err != nil {
				return fmt.Errorf("cannot unmarshal %s: %w", kv.Key, err)
			}
		case keyGithub:
			if err := json.Unmarshal(kv.Value, &s.Github); err != nil {
				return fmt.Errorf("cannot unmarshal %s: %w", kv.Key, err)
			}
		default:
			return fmt.Errorf("unknown key: %s", kv.Key)
		}
	}

	return nil
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
			Enabled: s.Github.Enabled,
			Token:   s.Github.Token,
		},
	}
}

func ToSvcSettings(s Settings) svcmodel.Settings {
	return svcmodel.Settings{
		OrganizationName:      s.OrganizationName,
		DefaultReleaseMessage: s.DefaultReleaseMessage,
		Slack: svcmodel.SlackSettings{
			Enabled: s.Slack.Enabled,
			Token:   s.Slack.Token,
		},
		Github: svcmodel.GithubSettings{
			Enabled: s.Github.Enabled,
			Token:   s.Github.Token,
		},
	}
}
