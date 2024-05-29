package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"release-manager/github/model"
	svcerrors "release-manager/service/errors"
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

func (c *Client) GetGithubRepoFromRawURL(rawURL string) (svcmodel.GithubRepo, error) {
	repo, err := model.ToSvcGithubRepo(rawURL)
	if err != nil {
		return svcmodel.GithubRepo{}, svcerrors.NewGithubRepositoryInvalidURL().Wrap(err).WithMessage(err.Error())
	}

	return repo, nil
}

func (c *Client) ReadTagsForRepository(ctx context.Context, tkn string, repoURL url.URL) ([]svcmodel.GitTag, error) {
	repo, err := model.ToGithubRepo(repoURL)
	if err != nil {
		return nil, svcerrors.NewGithubRepositoryInvalidURL().Wrap(err).WithMessage(err.Error())
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
				return nil, svcerrors.NewGithubClientUnauthorizedError().Wrap(err)
			case http.StatusForbidden:
				return nil, svcerrors.NewGithubClientForbiddenError().Wrap(err)
			case http.StatusNotFound:
				return nil, svcerrors.NewGithubRepositoryNotFoundError().Wrap(err)
			}
		}

		return nil, err
	}

	return model.ToSvcGitTags(t), nil
}

func (c *Client) ReadTagByName(ctx context.Context, tkn string, repoURL url.URL, tagName string) (svcmodel.GitTag, error) {
	repo, err := model.ToGithubRepo(repoURL)
	if err != nil {
		return svcmodel.GitTag{}, svcerrors.NewGithubRepositoryInvalidURL().Wrap(err).WithMessage(err.Error())
	}

	// Git tag can be fetched only by its SHA, using GET /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Another limitation is that only annotated tags can be fetched by /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Because lightweight tags do not have their own SHA, they only reference a commit
	// Docs https://docs.github.com/rest/git/tags#get-a-tag
	//
	// So in order to check if a tag exists by name (both lightweights and annotated tags), GET /repos/{owner}/{repo}/git/ref/{ref} is used
	// Docs https://docs.github.com/rest/git/refs#get-a-reference
	_, _, err = c.getGithubClient(tkn).Git.GetRef(
		ctx,
		repo.OwnerSlug,
		repo.RepositorySlug,
		fmt.Sprintf("tags/%s", tagName),
	)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Response.StatusCode == http.StatusNotFound {
			return svcmodel.GitTag{}, svcerrors.NewGitTagNotFoundError().Wrap(err)
		}

		return svcmodel.GitTag{}, err
	}

	return svcmodel.GitTag{Name: tagName}, nil
}

func (c *Client) getGithubClient(tkn string) *github.Client {
	return github.NewClient(nil).WithAuthToken(tkn)
}
