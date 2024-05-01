package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type GithubClient struct {
	mock.Mock
}

func (m *GithubClient) SetToken(token string) {
	m.Called(token)
}

func (m *GithubClient) ListTagsForRepository(ctx context.Context, repo svcmodel.GithubRepository) ([]svcmodel.GitTag, error) {
	args := m.Called(ctx, repo)
	return args.Get(0).([]svcmodel.GitTag), args.Error(1)
}
