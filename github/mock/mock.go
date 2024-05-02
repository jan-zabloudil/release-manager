package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (m *Client) RefreshClientWithToken(token string) {
	m.Called(token)
}

func (m *Client) ListTagsForRepository(ctx context.Context, repo svcmodel.GithubRepository) ([]svcmodel.GitTag, error) {
	args := m.Called(ctx, repo)
	return args.Get(0).([]svcmodel.GitTag), args.Error(1)
}

func (m *Client) CreateReleaseDraft(ctx context.Context, repo svcmodel.GithubRepository, input svcmodel.CreateReleaseDraftInput) (svcmodel.ReleaseDraft, error) {
	args := m.Called(ctx, repo, input)
	return args.Get(0).(svcmodel.ReleaseDraft), args.Error(1)
}
