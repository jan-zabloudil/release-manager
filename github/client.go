package github

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"release-manager/github/model"
	"release-manager/pkg/apierrors"
	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

const (
	tagsToFetch = 100
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ReadTagsForRepository(ctx context.Context, tkn string, repoURL url.URL) ([]svcmodel.GitTag, error) {
	repo, err := model.ToGithubRepo(repoURL)
	if err != nil {
		return nil, apierrors.NewGithubRepositoryInvalidURL().Wrap(err).WithMessage(err.Error())
	}

	// GitHub client is created in the function because it needs to be authenticated with a token (that is passed as an argument)
	client := github.NewClient(nil).WithAuthToken(tkn)

	// Up to 100 tags can be fetched per page
	// If the number of pages is not specified in the list options, only one page will be fetched
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-tags
	t, _, err := client.Repositories.ListTags(
		ctx,
		repo.OwnerSlug,
		repo.RepositorySlug,
		&github.ListOptions{PerPage: tagsToFetch},
	)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) {
			switch githubErr.Response.StatusCode {
			case http.StatusUnauthorized:
				return nil, apierrors.NewGithubClientUnauthorizedError().Wrap(err)
			case http.StatusForbidden:
				return nil, apierrors.NewGithubClientForbiddenError().Wrap(err)
			case http.StatusNotFound:
				return nil, apierrors.NewGithubRepositoryNotFoundError().Wrap(err)
			}
		}

		return nil, err
	}

	return model.ToSvcGitTags(t), nil
}
