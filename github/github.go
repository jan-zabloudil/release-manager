package github

import (
	"context"
	"fmt"

	"release-manager/github/model"
	"release-manager/github/util"
	"release-manager/pkg/githuberrors"
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

func (c *Client) RefreshClientWithToken(token string) {
	c.client = github.NewClient(nil).WithAuthToken(token)
}

func (c *Client) ListTagsForRepository(ctx context.Context, repo svcmodel.GithubRepository) ([]svcmodel.GitTag, error) {
	// Up to 100 tags can be fetched per page
	// If the number of pages is not specified in the list options, only one page will be fetched
	// Docs https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-tags
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

func (c *Client) CreateReleaseDraft(ctx context.Context, repo svcmodel.GithubRepository, input svcmodel.CreateReleaseDraftInput) (svcmodel.ReleaseDraft, error) {
	n, err := c.generateReleaseNotes(ctx, repo, input)
	if err != nil {
		return svcmodel.ReleaseDraft{}, fmt.Errorf("failed to generate release notes: %w", err)
	}
	draft := svcmodel.NewDraftRelease(n.ReleaseName, n.ReleaseNotes)

	t, err := c.getTagByName(ctx, repo, input.TagName)
	if err != nil && !githuberrors.IsNotFoundError(err) {
		return svcmodel.ReleaseDraft{}, fmt.Errorf("failed to get tag by name: %w", err)
	}
	if githuberrors.IsNotFoundError(err) {
		draft.LinkSourceCodeByNewTag(
			model.ToSvcGitTagInput(input.TagName, input.TargetCommitish),
		)

		// If the tag does not exist yet, there is no need to check if GitHub release exists
		return draft, nil
	}
	draft.LinkSourceCodeByExistingTag(t)

	r, err := c.getReleaseByTag(ctx, repo, input.TagName)
	if err != nil && !githuberrors.IsNotFoundError(err) {
		return svcmodel.ReleaseDraft{}, fmt.Errorf("failed to get release by tag: %w", err)
	}
	if githuberrors.IsNotFoundError(err) {
		return draft, nil
	}
	draft.AddGithubRelease(r)

	return draft, nil
}

func (c *Client) getTagByName(ctx context.Context, repo svcmodel.GithubRepository, tagName string) (svcmodel.GitTag, error) {
	// Git tag can be fetched only by its SHA, using GET /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Another limitation is that only annotated tags can be fetched by /repos/{owner}/{repo}/git/tags/{tag_sha}
	// Because lightweight tags do not have their own SHA, they only reference a commit
	// Docs https://docs.github.com/rest/git/tags#get-a-tag
	//
	// So in order to check if a tag exists by name (both lightweights and annotated tags), GET /repos/{owner}/{repo}/git/ref/{ref} is used
	// Docs https://docs.github.com/rest/git/refs#get-a-reference
	_, _, err := c.client.Git.GetRef(
		ctx,
		repo.OwnerSlug,
		repo.RepositorySlug,
		fmt.Sprintf("tags/%s", tagName),
	)
	if err != nil {
		return svcmodel.GitTag{}, util.ToGithubError(err)
	}

	return model.ToSvcGitTag(tagName), nil
}

func (c *Client) getReleaseByTag(ctx context.Context, repo svcmodel.GithubRepository, tagName string) (svcmodel.GithubRelease, error) {
	gr, _, err := c.client.Repositories.GetReleaseByTag(
		ctx,
		repo.OwnerSlug,
		repo.RepositorySlug,
		tagName,
	)
	if err != nil {
		return svcmodel.GithubRelease{}, util.ToGithubError(err)
	}

	r, err := model.ToSvcGithubRelease(gr)
	if err != nil {
		return svcmodel.GithubRelease{}, githuberrors.NewToSvcModelError().Wrap(err)
	}

	return r, nil
}

func (c *Client) generateReleaseNotes(ctx context.Context, repo svcmodel.GithubRepository, input svcmodel.CreateReleaseDraftInput) (model.GeneratedNotes, error) {
	// When creating a GitHub release (or generating release notes), the tag name is required
	// If tag does not exist yet, the release notes will be generated from the target commitish (can be a branch or a commit SHA)
	// If target commitish is not provided, the default branch will be used
	// Docs https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#generate-release-notes-content-for-a-release
	n, _, err := c.client.Repositories.GenerateReleaseNotes(ctx, repo.OwnerSlug, repo.RepositorySlug, &github.GenerateNotesOptions{
		TagName:         input.TagName,
		TargetCommitish: input.TargetCommitish,
		PreviousTagName: input.PreviousTagName,
	})
	if err != nil {
		return model.GeneratedNotes{}, util.ToGithubError(err) // TODO handle more specific errors such as invalid target commitish etc
	}

	return model.GeneratedNotes{
		ReleaseName:  n.Name,
		ReleaseNotes: n.Body,
	}, nil
}
