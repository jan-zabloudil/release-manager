package model

import (
	"fmt"
	"net/url"

	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

func ToSvcGitTags(tags []*github.RepositoryTag) []svcmodel.GitTag {
	t := make([]svcmodel.GitTag, 0, len(tags))
	for _, tag := range tags {
		t = append(t, svcmodel.GitTag{Name: *tag.Name})
	}

	return t
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
