package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"release-manager/github/model"
	"release-manager/github/util"
	svcerrors "release-manager/service/errors"
	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

const (
	tagsToFetch = 100

	// GitHub API error codes
	errCodeAlreadyExists = "already_exists"

	// GitHub API error fields
	gitTagNameField = "tag_name"
)

var (
	errGithubReleaseNotFound      = errors.New("github release not found")
	errGithubReleaseAlreadyExists = errors.New("github release already exists")
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ReadRepo(ctx context.Context, tkn string, rawRepoURL string) (svcmodel.GithubRepo, error) {
	ownerSlug, repoSlug, err := model.ParseGithubRepoURL(rawRepoURL)
	if err != nil {
		return svcmodel.GithubRepo{}, svcerrors.NewGithubRepoInvalidURL().Wrap(err).WithMessage(err.Error())
	}

	// Docs: https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
	repo, _, err := c.getGithubClient(tkn).Repositories.Get(ctx, ownerSlug, repoSlug)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Response.StatusCode == http.StatusNotFound {
			return svcmodel.GithubRepo{}, svcerrors.NewGithubRepoNotFoundError().Wrap(err)
		}

		return svcmodel.GithubRepo{}, util.TranslateGithubAuthError(err)
	}

	u, err := url.Parse(repo.GetHTMLURL())
	if err != nil {
		return svcmodel.GithubRepo{}, fmt.Errorf("failed to parse repo URL: %w", err)
	}

	return svcmodel.GithubRepo{
		OwnerSlug: ownerSlug,
		RepoSlug:  repoSlug,
		URL:       *u,
	}, nil
}

func (c *Client) ReadTagsForRepo(ctx context.Context, tkn string, repo svcmodel.GithubRepo) ([]svcmodel.GitTag, error) {
	// Up to 100 tags can be fetched per page
	// If the number of pages is not specified in the list options, only one page will be fetched
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-tags
	t, _, err := c.getGithubClient(tkn).Repositories.ListTags(
		ctx,
		repo.OwnerSlug,
		repo.RepoSlug,
		&github.ListOptions{PerPage: tagsToFetch},
	)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Response.StatusCode == http.StatusNotFound {
			return nil, svcerrors.NewGithubRepoNotFoundError().Wrap(err)
		}

		return nil, util.TranslateGithubAuthError(err)
	}

	return model.ToSvcGitTags(t), nil
}

func (c *Client) TagExists(ctx context.Context, tkn string, repo svcmodel.GithubRepo, tagName string) (bool, error) {
	// Git tag can be fetched only by its SHA, using GET /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Another limitation is that only annotated tags can be fetched by /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Because lightweight tags do not have their own SHA, they only reference a commit
	// Docs https://docs.github.com/rest/git/tags#get-a-tag
	//
	// So in order to check if a tag exists by name (both lightweights and annotated tags), GET /repos/{owner}/{repo}/git/ref/{ref} is used
	// Docs https://docs.github.com/rest/git/refs#get-a-reference
	_, _, err := c.getGithubClient(tkn).Git.GetRef(
		ctx,
		repo.OwnerSlug,
		repo.RepoSlug,
		fmt.Sprintf("tags/%s", tagName),
	)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Response.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, util.TranslateGithubAuthError(err)
	}

	return true, nil
}

func (c *Client) UpsertRelease(ctx context.Context, tkn string, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	if err := c.createRelease(ctx, tkn, repo, rls); err != nil {
		if errors.Is(err, errGithubReleaseAlreadyExists) {
			if err := c.updateRelease(ctx, tkn, repo, rls); err != nil {
				return fmt.Errorf("failed to update release: %w", err)
			}

			return nil
		}

		return fmt.Errorf("failed to create release: %w", err)
	}

	return nil
}

func (c *Client) DeleteReleaseByTag(ctx context.Context, tkn string, repo svcmodel.GithubRepo, tagName string) error {
	// Release can be deleted only by release ID
	// Therefore I need to get release ID first
	id, err := c.getReleaseIDByTag(ctx, tkn, repo, tagName)
	if err != nil {
		if errors.Is(err, errGithubReleaseNotFound) {
			return svcerrors.NewGithubReleaseNotFoundError().Wrap(err)
		}

		return fmt.Errorf("failed to get release ID: %w", err)
	}

	_, err = c.getGithubClient(tkn).Repositories.DeleteRelease(
		ctx,
		repo.OwnerSlug,
		repo.RepoSlug,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete release: %w", util.TranslateGithubAuthError(err))
	}

	return nil
}

// createRelease is an internal method for creating a release.
// returns internal errGithubReleaseAlreadyExists if the release already exists
func (c *Client) createRelease(ctx context.Context, tkn string, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	// Creates a new release
	// Docs: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#create-a-release
	//
	// TagName is the name of the tag to link the release to
	// Name is the name of the release
	// Body is the description of the release
	_, _, err := c.getGithubClient(tkn).Repositories.CreateRelease(ctx, repo.OwnerSlug, repo.RepoSlug, &github.RepositoryRelease{
		TagName: &rls.GitTagName,
		Name:    &rls.ReleaseTitle,
		Body:    &rls.ReleaseNotes,
	})
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Errors != nil {
			// GitHub returns error response as an array of errors
			// Each error contains fields (code, resource, field)
			for _, e := range githubErr.Errors {
				if e.Code == errCodeAlreadyExists && e.Field == gitTagNameField {
					return errGithubReleaseAlreadyExists
				}
			}
		}

		return util.TranslateGithubAuthError(err)
	}

	return nil
}

func (c *Client) updateRelease(ctx context.Context, tkn string, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	// Release can be updated only by release ID
	// Therefore I need to get release ID first
	id, err := c.getReleaseIDByTag(ctx, tkn, repo, rls.GitTagName)
	if err != nil {
		return fmt.Errorf("failed to get release ID: %w", err)
	}

	// Updates a release
	// Docs: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#update-a-release
	//
	// Name is the name of the release
	// Body is the description of the release
	_, _, err = c.getGithubClient(tkn).Repositories.EditRelease(
		ctx,
		repo.OwnerSlug,
		repo.RepoSlug,
		id,
		&github.RepositoryRelease{
			Name: &rls.ReleaseTitle,
			Body: &rls.ReleaseNotes,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update release: %w", util.TranslateGithubAuthError(err))
	}

	return nil
}

func (c *Client) GenerateRepoURL(ownerSlug, repoSlug string) (url.URL, error) {
	if ownerSlug == "" || repoSlug == "" {
		return url.URL{}, errors.New("empty owner or repo slug")
	}

	rawURL := fmt.Sprintf("https://github.com/%s/%s", ownerSlug, repoSlug)
	u, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, err
	}

	return *u, nil
}

func (c *Client) GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error) {
	if tagName == "" || ownerSlug == "" || repoSlug == "" {
		return url.URL{}, errors.New("empty tag name, owner or repo slug")
	}

	rawURL := fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", ownerSlug, repoSlug, tagName)
	u, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, err
	}

	return *u, nil
}

// getReleaseByTag is an internal method for fetching a release ID.
// This method simplifies the logic in other functions that need to get a release,
// as it also returns the private error errGithubReleaseNotFound if the release is not found.
// Other functions can then check if the error is equal to errGithubReleaseNotFound
// and handle the error based on their use case.
func (c *Client) getReleaseIDByTag(ctx context.Context, tkn string, repo svcmodel.GithubRepo, tagName string) (int64, error) {
	rls, _, err := c.getGithubClient(tkn).Repositories.GetReleaseByTag(
		ctx,
		repo.OwnerSlug,
		repo.RepoSlug,
		tagName,
	)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Response.StatusCode == http.StatusNotFound {
			return 0, errGithubReleaseNotFound
		}

		return 0, util.TranslateGithubAuthError(err)
	}

	return rls.GetID(), nil
}

func (c *Client) getGithubClient(tkn string) *github.Client {
	return github.NewClient(nil).WithAuthToken(tkn)
}
