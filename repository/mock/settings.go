package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type SettingsRepository struct {
	mock.Mock
}

func (m *SettingsRepository) Update(ctx context.Context, fn svcmodel.UpdateSettingsFunc) (svcmodel.Settings, error) {
	args := m.Called(ctx, fn)
	return args.Get(0).(svcmodel.Settings), args.Error(1)
}

func (m *SettingsRepository) Read(ctx context.Context) (svcmodel.Settings, error) {
	args := m.Called(ctx)
	return args.Get(0).(svcmodel.Settings), args.Error(1)
}
