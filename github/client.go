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

	// Up to 100 tags can be fetched per page
	// If the number of pages is not specified in the list options, only one page will be fetched
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-tags
	t, _, err := c.getGithubClient(tkn).Repositories.ListTags(
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

func (c *Client) CreateRelease(
	ctx context.Context,
	tkn string,
	repoURL url.URL,
	input svcmodel.CreateReleaseInput,
) (svcmodel.GithubRelease, error) {
	repo, err := model.ToGithubRepo(repoURL)
	if err != nil {
		return svcmodel.GithubRelease{}, apierrors.NewGithubRepositoryInvalidURL().Wrap(err).WithMessage(err.Error())
	}

	// Creates a new release
	// Docs: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#create-a-release
	//
	// TagName is the name of the tag to link the release to
	// Name is the name of the release
	// Body is the description of the release
	rls, _, err := c.getGithubClient(tkn).Repositories.CreateRelease(ctx, repo.OwnerSlug, repo.RepositorySlug, &github.RepositoryRelease{
		TagName: &input.GitTagName,
		Name:    &input.ReleaseTitle,
		Body:    &input.ReleaseNotes,
	})
	if err != nil {
		// TODO translate to service error if this function is not executed asynchronously
		return svcmodel.GithubRelease{}, err
	}

	return model.ToSvcGithubRelease(rls)
}

func (c *Client) getGithubClient(tkn string) *github.Client {
	return github.NewClient(nil).WithAuthToken(tkn)
}
