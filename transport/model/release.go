package model

import (
	"time"

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
	ID           uuid.UUID `json:"id"`
	ReleaseTitle string    `json:"release_title"`
	ReleaseNotes string    `json:"release_notes"`
	GitTagName   string    `json:"git_tag_name"`
	GitTagURL    string    `json:"git_tag_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
		ReleaseTitle: r.ReleaseTitle,
		ReleaseNotes: r.ReleaseNotes,
		GitTagName:   r.GitTagName,
		GitTagURL:    r.GitTagURL.String(),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func ToReleases(releases []svcmodel.Release) []Release {
	r := make([]Release, 0, len(releases))
	for _, release := range releases {
		r = append(r, ToRelease(release))
	}
	return r
}

type GithubGeneratedReleaseNotesInput struct {
	GitTagName         *string `json:"git_tag_name"`
	PreviousGitTagName *string `json:"previous_git_tag_name"`
}

type GithubGeneratedReleaseNotes struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
}

func ToSvcGithubGeneratedReleaseNotesInput(n GithubGeneratedReleaseNotesInput) svcmodel.GithubGeneratedReleaseNotesInput {
	return svcmodel.GithubGeneratedReleaseNotesInput{
		GitTagName:         n.GitTagName,
		PreviousGitTagName: n.PreviousGitTagName,
	}
}

func ToGithubGeneratedReleaseNotes(n svcmodel.GithubGeneratedReleaseNotes) GithubGeneratedReleaseNotes {
	return GithubGeneratedReleaseNotes{
		Title: n.Title,
		Notes: n.Notes,
	}
}
