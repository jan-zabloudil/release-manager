package model

import (
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
