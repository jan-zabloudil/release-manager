package model

import svcmodel "release-manager/service/model"

type GitTag struct {
	Name string `json:"name"`
}

func ToGitTags(tags []svcmodel.GitTag) []GitTag {
	t := make([]GitTag, 0, len(tags))
	for _, tag := range tags {
		t = append(t, GitTag{Name: tag.Name})
	}

	return t
}
