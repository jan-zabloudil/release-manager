package github

import (
	"context"

	"release-manager/github/model"
	"release-manager/github/util"
	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

const (
	tagsToFetch = 100
)

type Client struct {
	client *github.Client
}

func NewClient() *Client {
	return &Client{
		client: github.NewClient(nil),
	}
}

func (c *Client) SetToken(token string) {
	c.client = c.client.WithAuthToken(token)
}

func (c *Client) ListTagsForRepository(ctx context.Context, repo svcmodel.GithubRepository) ([]svcmodel.GitTag, error) {
	// Up to 100 tags can be fetched per page
	// If the number of pages is not specified in the list options, only one page will be fetched
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-tags
	t, _, err := c.client.Repositories.ListTags(
		ctx,
		repo.OwnerSlug,
		repo.RepositorySlug,
		&github.ListOptions{PerPage: tagsToFetch},
	)
	if err != nil {
		return nil, util.ToGithubError(err)
	}

	return model.ToSvcGitTags(t), nil
}
