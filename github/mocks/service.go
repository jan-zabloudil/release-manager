package mocks

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type MockGitHubService struct {
	mock.Mock
}

func (m *MockGitHubService) ListTags(ctx context.Context, owner, repo string) ([]svcmodel.GitTag, error) {
	args := m.Called(ctx, owner, repo)
	return args.Get(0).([]svcmodel.GitTag), args.Error(1)
}
