package model

import (
	"net/url"

	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

type GeneratedNotes struct {
	ReleaseName  string
	ReleaseNotes string
}

func ToSvcGitTags(tags []*github.RepositoryTag) []svcmodel.GitTag {
	t := make([]svcmodel.GitTag, 0, len(tags))
	for _, tag := range tags {
		if tag.Name != nil {
			t = append(t, ToSvcGitTag(*tag.Name))
		}
	}

	return t
}

func ToSvcGitTag(name string) svcmodel.GitTag {
	return svcmodel.GitTag{Name: name}
}

func ToSvcGithubRelease(r *github.RepositoryRelease) (svcmodel.GithubRelease, error) {
	var svcRelease svcmodel.GithubRelease

	if r.HTMLURL != nil {
		u, err := url.Parse(*r.HTMLURL)
		if err != nil {
			return svcmodel.GithubRelease{}, err
		}

		svcRelease.ReleasePageURL = *u
	}

	if r.Name != nil {
		svcRelease.Name = *r.Name
	}

	if r.Body != nil {
		svcRelease.ReleaseNotes = *r.Body
	}

	return svcRelease, nil
}

func ToSvcGitTagInput(tagName string, targetCommitish *string) svcmodel.GitTagInput {
	if targetCommitish == nil {
		return svcmodel.GitTagInput{
			TagName: tagName,
		}
	}

	return svcmodel.GitTagInput{
		TagName:         tagName,
		TargetCommitish: *targetCommitish,
	}
}
