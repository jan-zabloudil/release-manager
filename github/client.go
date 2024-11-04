package github

import (
	"context"
	"fmt"
	"net/url"

	"release-manager/github/model"
	"release-manager/github/util"
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

func (c *Client) ReadRepo(ctx context.Context, tkn svcmodel.GithubToken, rawRepoURL string) (svcmodel.GithubRepo, error) {
	ownerSlug, repoSlug, err := util.ParseGithubRepoURL(rawRepoURL)
	if err != nil {
		return svcmodel.GithubRepo{}, svcerrors.NewGithubRepoInvalidURL().Wrap(err).WithMessage(err.Error())
	}

	return withGithubClientResult[svcmodel.GithubRepo](tkn, func(client *github.Client) (svcmodel.GithubRepo, error) {
		// Docs: https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
		repo, _, err := client.Repositories.Get(ctx, ownerSlug, repoSlug)
		if err != nil {
			if util.IsNotFoundError(err) {
				return svcmodel.GithubRepo{}, svcerrors.NewGithubRepoNotFoundError().Wrap(err)
			}

			return svcmodel.GithubRepo{}, err
		}

		return model.ToSvcGithubRepo(repo, ownerSlug, repoSlug)
	})
}

func (c *Client) ReadTagsForRepo(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo) ([]svcmodel.GitTag, error) {
	return withGithubClientResult[[]svcmodel.GitTag](tkn, func(client *github.Client) ([]svcmodel.GitTag, error) {
		// Up to 100 tags can be fetched per page
		// If the number of pages is not specified in the list options, only one page will be fetched
		// Docs: https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-tags
		t, _, err := client.Repositories.ListTags(
			ctx,
			repo.OwnerSlug,
			repo.RepoSlug,
			&github.ListOptions{PerPage: tagsToFetch},
		)
		if err != nil {
			if util.IsNotFoundError(err) {
				return nil, svcerrors.NewGithubRepoNotFoundError().Wrap(err)
			}

			return nil, err
		}

		return model.ToSvcGitTags(t), nil
	})
}

func (c *Client) TagExists(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, tagName string) (bool, error) {
	// Git tag can be fetched only by its SHA, using GET /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Another limitation is that only annotated tags can be fetched by /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Because lightweight tags do not have their own SHA, they only reference a commit
	// Docs https://docs.github.com/rest/git/tags#get-a-tag
	//
	// So in order to check if a tag exists by name (both lightweights and annotated tags), GET /repos/{owner}/{repo}/git/ref/{ref} is used
	// Docs https://docs.github.com/rest/git/refs#get-a-reference
	err := withGithubClient(tkn, func(client *github.Client) error {
		_, _, err := client.Git.GetRef(
			ctx,
			repo.OwnerSlug,
			repo.RepoSlug,
			fmt.Sprintf("tags/%s", tagName),
		)
		return err
	})
	if err != nil {
		if util.IsNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *Client) UpsertRelease(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	if err := c.createRelease(ctx, tkn, repo, rls); err != nil {
		if util.IsReleaseAlreadyExistsError(err) {
			if err := c.updateRelease(ctx, tkn, repo, rls); err != nil {
				return fmt.Errorf("updating release: %w", err)
			}

			return nil
		}

		return fmt.Errorf("creating release: %w", err)
	}

	return nil
}

func (c *Client) DeleteReleaseByTag(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, tagName string) error {
	return withGithubClient(tkn, func(client *github.Client) error {
		// Release can be deleted only by release ID
		// Therefore I need to get release object first
		rls, _, err := client.Repositories.GetReleaseByTag(
			ctx,
			repo.OwnerSlug,
			repo.RepoSlug,
			tagName,
		)
		if err != nil {
			if util.IsNotFoundError(err) {
				return svcerrors.NewGithubReleaseNotFoundError().Wrap(err)
			}

			return fmt.Errorf("getting release by tag: %w", err)
		}

		if _, err := client.Repositories.DeleteRelease(
			ctx,
			repo.OwnerSlug,
			repo.RepoSlug,
			rls.GetID(),
		); err != nil {
			return fmt.Errorf("deleting release: %w", err)
		}

		return nil
	})
}

func (c *Client) GenerateReleaseNotes(
	ctx context.Context,
	tkn svcmodel.GithubToken,
	repo svcmodel.GithubRepo,
	input svcmodel.GithubReleaseNotesInput,
) (svcmodel.GithubReleaseNotes, error) {
	return withGithubClientResult[svcmodel.GithubReleaseNotes](tkn, func(client *github.Client) (svcmodel.GithubReleaseNotes, error) {
		// Generates release notes based on git tag and previous git tag
		// Git tag must be present, and it can be either existing tag or new tag that will be created
		// Previous git tag name is optional field
		//
		// Docs: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#generate-release-notes-content-for-a-release
		notes, _, err := client.Repositories.GenerateReleaseNotes(
			ctx,
			repo.OwnerSlug,
			repo.RepoSlug,
			&github.GenerateNotesOptions{
				TagName:         input.GetGitTagName(),
				PreviousTagName: input.PreviousGitTagName,
			},
		)
		if err != nil {
			if util.IsInvalidPreviousTagError(err) {
				return svcmodel.GithubReleaseNotes{},
					svcerrors.NewGithubNotesInvalidInputError().Wrap(err).WithMessage("Invalid previous git tag")
			}

			return svcmodel.GithubReleaseNotes{}, fmt.Errorf("generating release notes: %w", err)
		}

		return model.ToGithubReleaseNotes(notes), nil
	})
}

func (c *Client) GenerateRepoURL(ownerSlug, repoSlug string) (url.URL, error) {
	return util.GenerateRepoURL(ownerSlug, repoSlug)
}

func (c *Client) GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error) {
	return util.GenerateGitTagURL(ownerSlug, repoSlug, tagName)
}

func (c *Client) createRelease(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	return withGithubClient(tkn, func(client *github.Client) error {
		// Creates a new release
		// Docs: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#create-a-release
		//
		// TagName is the name of the tag to link the release to
		// Name is the name of the release
		// Body is the description of the release
		if _, _, err := client.Repositories.CreateRelease(ctx, repo.OwnerSlug, repo.RepoSlug, &github.RepositoryRelease{
			TagName: &rls.GitTagName,
			Name:    &rls.ReleaseTitle,
			Body:    &rls.ReleaseNotes,
		}); err != nil {
			return err
		}

		return nil
	})
}

func (c *Client) updateRelease(ctx context.Context, tkn svcmodel.GithubToken, repo svcmodel.GithubRepo, rls svcmodel.Release) error {
	return withGithubClient(tkn, func(client *github.Client) error {
		// Release can be updated only by release ID
		// Therefore I need to get release ID first
		githubRls, _, err := client.Repositories.GetReleaseByTag(
			ctx,
			repo.OwnerSlug,
			repo.RepoSlug,
			rls.GitTagName,
		)
		if err != nil {
			return fmt.Errorf("getting release by tag: %w", err)
		}

		// Updates a release
		// Docs: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#update-a-release
		//
		// Name is the name of the release
		// Body is the description of the release
		if _, _, err = client.Repositories.EditRelease(
			ctx,
			repo.OwnerSlug,
			repo.RepoSlug,
			githubRls.GetID(),
			&github.RepositoryRelease{
				Name: &rls.ReleaseTitle,
				Body: &rls.ReleaseNotes,
			},
		); err != nil {
			return fmt.Errorf("updating release: %w", err)
		}

		return nil
	})
}

func withGithubClientResult[T any](tkn svcmodel.GithubToken, fn func(client *github.Client) (T, error)) (T, error) {
	client := github.NewClient(nil).WithAuthToken(tkn.String())
	var zeroValue T
	result, err := fn(client)
	if err != nil {
		return zeroValue, util.TranslateGithubAuthError(err)
	}
	return result, nil
}

func withGithubClient(tkn svcmodel.GithubToken, fn func(client *github.Client) error) error {
	client := github.NewClient(nil).WithAuthToken(tkn.String())
	if err := fn(client); err != nil {
		return util.TranslateGithubAuthError(err)
	}
	return nil
}
