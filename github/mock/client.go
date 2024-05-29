package mock

import (
	"context"
	"net/url"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (c *Client) ReadTagsForRepo(ctx context.Context, tkn string, repoURL url.URL) ([]svcmodel.GitTag, error) {
	args := c.Called(ctx, tkn, repoURL)
	return args.Get(0).([]svcmodel.GitTag), args.Error(1)
}

func (c *Client) ReadRepo(ctx context.Context, tkn string, rawRepoURL string) (svcmodel.GithubRepo, error) {
	args := c.Called(ctx, tkn, rawRepoURL)
	return args.Get(0).(svcmodel.GithubRepo), args.Error(1)
}
