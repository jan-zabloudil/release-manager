package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type SettingsService struct {
	mock.Mock
}

func (m *SettingsService) GetGithubToken(ctx context.Context) (model.GithubToken, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.GithubToken), args.Error(1)
}

func (m *SettingsService) GetSlackToken(ctx context.Context) (model.SlackToken, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.SlackToken), args.Error(1)
}

func (m *SettingsService) GetDefaultReleaseMessage(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *SettingsService) GetGithubSettings(ctx context.Context) (model.GithubSettings, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.GithubSettings), args.Error(1)
}
