package model

import (
	"fmt"
	"net/url"

	"release-manager/github/util"
	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

type TagDeletionWebhookInput struct {
	Tag string `json:"ref" validate:"required"`
	// RefType is the type of the reference (e.g. "branch" or "tag")
	// We only care about tags
	RefType string `json:"ref_type" validate:"required,eq=tag"`
	Repo    struct {
		// Owner and repo slug of the GitHub repo separated by a slash
		// (e.g. "owner/repo")
		Slugs string `json:"full_name" validate:"required"`
	} `json:"repository"`
}

func ToSvcGitTag(tagName string, repo svcmodel.GithubRepo) (svcmodel.GitTag, error) {
	tagURL, err := util.GenerateGitTagURL(repo.OwnerSlug, repo.RepoSlug, tagName)
	if err != nil {
		return svcmodel.GitTag{}, err
	}

	return svcmodel.GitTag{
		Name: tagName,
		URL:  tagURL,
	}, nil
}

func ToSvcGitTags(tags []*github.RepositoryTag, repo svcmodel.GithubRepo) ([]svcmodel.GitTag, error) {
	t := make([]svcmodel.GitTag, 0, len(tags))
	for _, tag := range tags {
		svcTag, err := ToSvcGitTag(*tag.Name, repo)
		if err != nil {
			return nil, err
		}

		t = append(t, svcTag)
	}

	return t, nil
}

func ToSvcGithubRepo(repo *github.Repository, ownerSlug, repoSlug string) (svcmodel.GithubRepo, error) {
	u, err := url.Parse(repo.GetHTMLURL())
	if err != nil {
		return svcmodel.GithubRepo{}, fmt.Errorf("parsing GitHub repo URL: %w", err)
	}

	return svcmodel.GithubRepo{
		OwnerSlug: ownerSlug,
		RepoSlug:  repoSlug,
		URL:       *u,
	}, nil
}

func ToGithubReleaseNotes(notes *github.RepositoryReleaseNotes) svcmodel.GithubReleaseNotes {
	return svcmodel.GithubReleaseNotes{
		Title: notes.Name,
		Notes: notes.Body,
	}
}
