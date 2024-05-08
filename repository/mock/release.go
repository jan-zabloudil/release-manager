package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type ReleaseRepository struct {
	mock.Mock
}

func (m *ReleaseRepository) Create(ctx context.Context, rls svcmodel.Release) error {
	args := m.Called(ctx, rls)
	return args.Error(0)
}
