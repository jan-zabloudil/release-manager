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

func (c *Client) ReadTagsForRepo(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo) ([]svcmodel.GitTag, error) {
	args := c.Called(ctx, tkn, repo)
	return args.Get(0).([]svcmodel.GitTag), args.Error(1)
}

func (c *Client) ReadRepo(ctx context.Context, tkn svcmodel.GithubToken, rawRepoURL string) (svcmodel.GithubRepo, error) {
	args := c.Called(ctx, tkn, rawRepoURL)
	return args.Get(0).(svcmodel.GithubRepo), args.Error(1)
}

func (c *Client) GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error) {
	args := c.Called(ownerSlug, repoSlug, tagName)
	return args.Get(0).(url.URL), args.Error(1)
}

func (c *Client) DeleteReleaseByTag(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, tag svcmodel.GitTag) error {
	args := c.Called(ctx, tkn, repo, tag)
	return args.Error(0)
}

func (c *Client) ReadTag(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, tagName string) (svcmodel.GitTag, error) {
	args := c.Called(ctx, tkn, repo, tagName)
	return args.Get(0).(svcmodel.GitTag), args.Error(1)
}

func (c *Client) UpsertRelease(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	args := c.Called(ctx, tkn, repo, rls)
	return args.Error(0)
}

func (c *Client) GenerateRepoURL(ownerSlug, repoSlug string) (url.URL, error) {
	args := c.Called(ownerSlug, repoSlug)
	return args.Get(0).(url.URL), args.Error(1)
}

func (c *Client) GenerateReleaseNotes(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, input svcmodel.GithubReleaseNotesInput) (svcmodel.GithubReleaseNotes, error) {
	args := c.Called(ctx, tkn, repo, input)
	return args.Get(0).(svcmodel.GithubReleaseNotes), args.Error(1)
}

func (c *Client) ParseTagDeletionWebhook(ctx context.Context, webhook svcmodel.GithubTagDeletionWebhookInput, tkn svcmodel.GithubToken, secret svcmodel.GithubWebhookSecret) (svcmodel.GithubTagDeletionWebhookOutput, error) {
	args := c.Called(ctx, webhook, tkn, secret)
	return args.Get(0).(svcmodel.GithubTagDeletionWebhookOutput), args.Error(1)
}
