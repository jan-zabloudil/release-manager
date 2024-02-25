package model

import (
	svcmodel "release-manager/service/model"
)

type ReleaseMessage struct {
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	Includes Includes `json:"includes"`
}

type Includes struct {
	ProjectName   bool `json:"project_name"`
	AppName       bool `json:"app_name"`
	ReleaseName   bool `json:"release_name"`
	Changelog     bool `json:"changelog"`
	Deployments   bool `json:"deployments"`
	GithubRelease bool `json:"github_release"`
	GithubTag     bool `json:"github_tag"`
}

func ToDBReleaseMsg(title, text string, incl svcmodel.Includes) ReleaseMessage {
	return ReleaseMessage{
		Title:    title,
		Text:     text,
		Includes: Includes(incl),
	}
}

func ToSvcReleaseMsg(title, text string, incl Includes) svcmodel.ReleaseMessage {
	return svcmodel.ReleaseMessage{
		Title:    title,
		Text:     text,
		Includes: svcmodel.Includes(incl),
	}
}
