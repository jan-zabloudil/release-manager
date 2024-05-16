package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SettingsService struct {
	mock.Mock
}

func (m *SettingsService) Get(ctx context.Context, authUserID uuid.UUID) (model.Settings, error) {
	args := m.Called(ctx, authUserID)
	return args.Get(0).(model.Settings), args.Error(1)
}

func (m *SettingsService) Update(ctx context.Context, u model.UpdateSettingsInput, authUserID uuid.UUID) (model.Settings, error) {
	args := m.Called(ctx, u, authUserID)
	return args.Get(0).(model.Settings), args.Error(1)
}

func (m *SettingsService) GetGithubToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.Get(0).(string), args.Error(1)
}

func (m *SettingsService) GetSlackToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
