package model

import (
	"fmt"
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

// GithubRepo holds the owner and repository slugs of a GitHub repository
// Example URL: https://github.com/owner/repo, OwnerSlug: owner, RepositorySlug: repo
// Both slugs are needed for the GitHub API
type GithubRepo struct {
	OwnerSlug      string
	RepositorySlug string
}

func ToGithubRepo(u url.URL) (GithubRepo, error) {
	path := strings.Trim(u.Path, "/")
	slugs := strings.Split(path, "/")

	if len(slugs) != expectedGithubRepositoryURLSlugCount {
		return GithubRepo{}, fmt.Errorf("invalid GitHub repository URL: %s", u.String())
	}

	if slugs[0] == "" || slugs[1] == "" {
		return GithubRepo{}, fmt.Errorf("invalid GitHub repository URL: %s", u.String())
	}

	return GithubRepo{
		OwnerSlug:      slugs[0],
		RepositorySlug: slugs[1],
	}, nil
}

func ToSvcGitTags(tags []*github.RepositoryTag) []svcmodel.GitTag {
	t := make([]svcmodel.GitTag, 0, len(tags))
	for _, tag := range tags {
		t = append(t, svcmodel.GitTag{Name: *tag.Name})
	}

	return t
}

func ToSvcGithubRelease(release *github.RepositoryRelease) (svcmodel.GithubRelease, error) {
	u, err := url.Parse(release.GetHTMLURL())
	if err != nil {
		return svcmodel.GithubRelease{}, err
	}

	return svcmodel.GithubRelease{
		ID:          release.GetID(),
		Name:        release.GetName(),
		Body:        release.GetBody(),
		HTMLURL:     *u,
		TagName:     release.GetTagName(),
		CreatedAt:   release.CreatedAt.Time,
		PublishedAt: release.PublishedAt.Time,
	}, nil
}
