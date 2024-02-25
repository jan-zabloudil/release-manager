package model

import (
	svcmodel "release-manager/service/model"
)

const (
	projectNameKey   = "project_name"
	appNameKey       = "app_name"
	releaseNameKey   = "release_name"
	changelogKey     = "changelog"
	deploymentsKey   = "deployments"
	githubReleaseKey = "github_release"
	githubTagKey     = "github_tag"
)

type ReleaseMessage struct {
	Title    *string   `json:"title" validate:"required"`
	Text     *string   `json:"text"`
	Includes *[]string `json:"includes" validate:"required"`
}

type ReleaseMessagePatch struct {
	Title    *string   `json:"title"`
	Text     *string   `json:"text"`
	Includes *[]string `json:"includes"`
}

func NewSvcReleaseMsg(title, text *string, incl *[]string) svcmodel.ReleaseMessage {
	return ToSvcReleaseMsg(svcmodel.ReleaseMessage{}, title, text, incl)
}

func ToSvcReleaseMsg(msg svcmodel.ReleaseMessage, title, text *string, incl *[]string) svcmodel.ReleaseMessage {
	if title != nil {
		msg.Title = *title
	}
	if text != nil {
		msg.Text = *text
	}
	if incl != nil {
		// without resetting the struct, true values that are not provided in the payload would not be overridden
		msg.Includes = svcmodel.Includes{}

		allowedIncludes := map[string]*bool{
			projectNameKey:   &msg.Includes.ProjectName,
			appNameKey:       &msg.Includes.AppName,
			releaseNameKey:   &msg.Includes.ReleaseName,
			changelogKey:     &msg.Includes.Changelog,
			deploymentsKey:   &msg.Includes.Deployments,
			githubReleaseKey: &msg.Includes.GithubRelease,
			githubTagKey:     &msg.Includes.GithubTag,
		}

		for _, include := range *incl {
			if ptr, ok := allowedIncludes[include]; ok {
				*ptr = true
			}
		}
	}

	return msg
}

func ToNetReleaseMsg(title, text string, incl svcmodel.Includes) ReleaseMessage {
	m := map[string]bool{
		projectNameKey:   incl.ProjectName,
		appNameKey:       incl.AppName,
		releaseNameKey:   incl.ReleaseName,
		changelogKey:     incl.Changelog,
		deploymentsKey:   incl.Deployments,
		githubReleaseKey: incl.GithubRelease,
		githubTagKey:     incl.GithubTag,
	}

	var includes []string
	for field, value := range m {
		if value {
			includes = append(includes, field)
		}
	}

	return ReleaseMessage{
		Title:    &title,
		Text:     &text,
		Includes: &includes,
	}
}

func ToNetReleaseMsgs(msgs []svcmodel.ReleaseMessage) []ReleaseMessage {
	m := make([]ReleaseMessage, 0, len(msgs))
	for _, message := range msgs {
		m = append(m, ToNetReleaseMsg(message.Title, message.Text, message.Includes))
	}

	return m
}
