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

func (c *Client) ReadTagsForRepository(ctx context.Context, tkn string, repoURL url.URL) ([]svcmodel.GitTag, error) {
	args := c.Called(ctx, tkn, repoURL)
	return args.Get(0).([]svcmodel.GitTag), args.Error(1)
}

func (c *Client) GenerateGitTagURL(repoURL url.URL, tagName string) (*url.URL, error) {
	args := c.Called(repoURL, tagName)
	return args.Get(0).(*url.URL), args.Error(1)
}
