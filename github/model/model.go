package model

import (
	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

func ToRepositoryRelease(tagName, targetCommitish, name, body string) github.RepositoryRelease {
	return github.RepositoryRelease{
		TagName:         github.String(tagName),
		TargetCommitish: github.String(targetCommitish),
		Name:            github.String(name),
		Body:            github.String(body),
		Draft:           github.Bool(false),
		Prerelease:      github.Bool(false),
	}
}

func ToSvcGithubRelease(tagName, targetCommitish, name, changelog *string) svcmodel.GithubRelease {
	return svcmodel.GithubRelease{
		TagName:         *tagName,
		TargetCommitish: *targetCommitish,
		Name:            *name,
		Changelog:       *changelog,
	}
}

func ToSvcGitTags(tags []*github.RepositoryTag) []svcmodel.GitTag {
	t := make([]svcmodel.GitTag, 0, len(tags))
	for _, tag := range tags {
		t = append(t, svcmodel.GitTag{Name: *tag.Name})
	}

	return t
}
