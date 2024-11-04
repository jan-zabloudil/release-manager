package model

import (
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateReleaseInput struct {
	ReleaseTitle string `json:"release_title"`
	ReleaseNotes string `json:"release_notes"`
	GitTagName   string `json:"git_tag_name"`
}

type UpdateReleaseInput struct {
	ReleaseTitle *string `json:"release_title"`
	ReleaseNotes *string `json:"release_notes"`
}

type DeleteReleaseInput struct {
	DeleteGithubRelease bool `json:"delete_github_release"`
}

type Release struct {
	ID           id.Release          `json:"id"`
	ProjectID    id.Project          `json:"project_id"`
	ReleaseTitle string              `json:"release_title"`
	ReleaseNotes string              `json:"release_notes"`
	GitTagName   string              `json:"git_tag_name"`
	GitTagURL    string              `json:"git_tag_url"`
	Attachments  []ReleaseAttachment `json:"attachments"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type ReleaseAttachment struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	URL  string    `json:"url"`
}

func ToSvcCreateReleaseInput(r CreateReleaseInput) svcmodel.CreateReleaseInput {
	return svcmodel.CreateReleaseInput{
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
		GitTagName:   r.GitTagName,
	}
}

func ToSvcUpdateReleaseInput(r UpdateReleaseInput) svcmodel.UpdateReleaseInput {
	return svcmodel.UpdateReleaseInput{
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
	}
}

func ToSvcDeleteReleaseInput(r DeleteReleaseInput) svcmodel.DeleteReleaseInput {
	return svcmodel.DeleteReleaseInput{
		DeleteGithubRelease: r.DeleteGithubRelease,
	}
}

func ToRelease(r svcmodel.Release) Release {
	return Release{
		ID:           r.ID,
		ProjectID:    r.ProjectID,
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
		GitTagName:   r.GitTagName,
		GitTagURL:    r.GitTagURL.String(),
		Attachments:  ToReleaseAttachments(r.Attachments),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func ToReleaseAttachments(attachments []svcmodel.ReleaseAttachment) []ReleaseAttachment {
	a := make([]ReleaseAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		a = append(a, ToReleaseAttachment(attachment))
	}
	return a
}

func ToReleaseAttachment(a svcmodel.ReleaseAttachment) ReleaseAttachment {
	return ReleaseAttachment{
		ID:   a.ID,
		Name: a.Name,
		URL:  a.URL.String(),
	}
}

func ToReleases(releases []svcmodel.Release) []Release {
	r := make([]Release, 0, len(releases))
	for _, release := range releases {
		r = append(r, ToRelease(release))
	}
	return r
}

type GithubReleaseNotesInput struct {
	GitTagName         *string `json:"git_tag_name"`
	PreviousGitTagName *string `json:"previous_git_tag_name"`
}

type GithubReleaseNotes struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
}

func ToSvcGithubReleaseNotesInput(n GithubReleaseNotesInput) svcmodel.GithubReleaseNotesInput {
	return svcmodel.GithubReleaseNotesInput{
		GitTagName:         n.GitTagName,
		PreviousGitTagName: n.PreviousGitTagName,
	}
}

func ToGithubReleaseNotes(n svcmodel.GithubReleaseNotes) GithubReleaseNotes {
	return GithubReleaseNotes{
		Title: n.Title,
		Notes: n.Notes,
	}
}
