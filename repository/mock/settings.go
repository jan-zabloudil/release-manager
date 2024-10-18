package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type SettingsRepository struct {
	mock.Mock
}

func (m *SettingsRepository) Upsert(
	ctx context.Context,
	upsertFn func(svcmodel.Settings) (svcmodel.Settings, error),
) error {
	args := m.Called(ctx, upsertFn)
	return args.Error(0)
}

func (m *SettingsRepository) Read(ctx context.Context) (svcmodel.Settings, error) {
	args := m.Called(ctx)
	return args.Get(0).(svcmodel.Settings), args.Error(1)
}
