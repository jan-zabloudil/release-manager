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

func (c *Client) ReadTagsForRepo(ctx context.Context, tkn string, repo svcmodel.GithubRepo) ([]svcmodel.GitTag, error) {
	args := c.Called(ctx, tkn, repo)
	return args.Get(0).([]svcmodel.GitTag), args.Error(1)
}

func (c *Client) ReadRepo(ctx context.Context, tkn string, rawRepoURL string) (svcmodel.GithubRepo, error) {
	args := c.Called(ctx, tkn, rawRepoURL)
	return args.Get(0).(svcmodel.GithubRepo), args.Error(1)
}

func (c *Client) GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error) {
	args := c.Called(ownerSlug, repoSlug, tagName)
	return args.Get(0).(url.URL), args.Error(1)
}

func (c *Client) DeleteReleaseByTag(ctx context.Context, tkn string, repo svcmodel.GithubRepo, tagName string) error {
	args := c.Called(ctx, tkn, repo, tagName)
	return args.Error(0)
}

func (c *Client) TagExists(ctx context.Context, tkn string, repo svcmodel.GithubRepo, tagName string) (bool, error) {
	args := c.Called(ctx, tkn, repo, tagName)
	return args.Bool(0), args.Error(1)
}

func (c *Client) UpsertRelease(ctx context.Context, tkn string, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	args := c.Called(ctx, tkn, repo, rls)
	return args.Error(0)
}

func (c *Client) GenerateRepoURL(ownerSlug, repoSlug string) (url.URL, error) {
	args := c.Called(ownerSlug, repoSlug)
	return args.Get(0).(url.URL), args.Error(1)
}

func (c *Client) GenerateReleaseNotes(ctx context.Context, token string, repo svcmodel.GithubRepo, input svcmodel.GithubGeneratedReleaseNotesInput) (svcmodel.GithubGeneratedReleaseNotes, error) {
	args := c.Called(ctx, token, repo, input)
	return args.Get(0).(svcmodel.GithubGeneratedReleaseNotes), args.Error(1)
}
