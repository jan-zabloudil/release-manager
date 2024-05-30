package model

import (
	"errors"
	"net/url"
	"strings"

	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

const (
	// expectedGithubRepositoryURLSlugCount is the expected number of slugs in a GitHub repository URL
	// Example URL: https://github.com/owner/repo -> owner and repo are the slugs
	expectedGithubRepositoryURLSlugCount = 2
)

var (
	errInvalidGithubRepoURLPath = errors.New("invalid GitHub repository URL path, not in the format /owner/repo")
)

func ParseGithubRepoURL(rawURL string) (ownerSlug, repoSlug string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	// GitHub repo URL format: https://github.com/owner/repo,
	// OwnerSlug: owner, RepoSlug: repo.
	path := strings.Trim(u.Path, "/")
	slugs := strings.Split(path, "/")

	if len(slugs) != expectedGithubRepositoryURLSlugCount {
		return "", "", errInvalidGithubRepoURLPath
	}

	if slugs[0] == "" || slugs[1] == "" {
		return "", "", errInvalidGithubRepoURLPath
	}

	return slugs[0], slugs[1], nil
}

func ToSvcGitTags(tags []*github.RepositoryTag) []svcmodel.GitTag {
	t := make([]svcmodel.GitTag, 0, len(tags))
	for _, tag := range tags {
		t = append(t, svcmodel.GitTag{Name: *tag.Name})
	}

	return t
}
