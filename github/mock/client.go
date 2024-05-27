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

func (c *Client) ReadTagByName(ctx context.Context, tkn string, repoURL url.URL, tagName string) (svcmodel.GitTag, error) {
	args := c.Called(ctx, tkn, repoURL, tagName)
	return args.Get(0).(svcmodel.GitTag), args.Error(1)
}

func (c *Client) CreateRelease(ctx context.Context, tkn string, repoURL url.URL, input svcmodel.CreateReleaseInput) (svcmodel.GithubRelease, error) {
	args := c.Called(ctx, tkn, repoURL, input)
	return args.Get(0).(svcmodel.GithubRelease), args.Error(1)
}

func (c *Client) ReadReleaseByTag(ctx context.Context, tkn string, repoURL url.URL, tagName string) (svcmodel.GithubRelease, error) {
	args := c.Called(ctx, tkn, repoURL, tagName)
	return args.Get(0).(svcmodel.GithubRelease), args.Error(1)
}
