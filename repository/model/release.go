package model

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Release struct {
	ID           id.Release  `db:"id"`
	ProjectID    id.Project  `db:"project_id"`
	ReleaseTitle string      `db:"release_title"`
	ReleaseNotes string      `db:"release_notes"`
	AuthorUserID id.AuthUser `db:"created_by"`
	GitTagName   string      `db:"git_tag_name"`
	// GithubRepoSlug and GithubOwnerSlug are fetched from the project
	// and are used to generate the tag URL
	GithubRepoSlug  sql.NullString      `db:"github_repo_slug"`
	GithubOwnerSlug sql.NullString      `db:"github_owner_slug"`
	Attachments     []ReleaseAttachment `db:"attachments"`
	CreatedAt       time.Time           `db:"created_at"`
	UpdatedAt       time.Time           `db:"updated_at"`
}

type gitTagURLGeneratorFunc func(ownerSlug, repoSlug, tag string) (url.URL, error)

func ToSvcRelease(
	rls Release,
	tagURLGenerator gitTagURLGeneratorFunc,
	fileURLGenerator fileURLGeneratorFunc,
) (svcmodel.Release, error) {
	tagURL, err := tagURLGenerator(rls.GithubOwnerSlug.String, rls.GithubRepoSlug.String, rls.GitTagName)
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("generating a git tag URL: %w", err)
	}

	attachments, err := ToSvcReleaseAttachments(rls.Attachments, fileURLGenerator)
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("converting release attachments to service model: %w", err)
	}

	return svcmodel.Release{
		ID:           rls.ID,
		ProjectID:    rls.ProjectID,
		ReleaseTitle: rls.ReleaseTitle,
		ReleaseNotes: rls.ReleaseNotes,
		GitTagName:   rls.GitTagName,
		GitTagURL:    tagURL,
		AuthorUserID: rls.AuthorUserID,
		Attachments:  attachments,
		CreatedAt:    rls.CreatedAt,
		UpdatedAt:    rls.UpdatedAt,
	}, nil
}

func ToSvcReleases(
	releases []Release,
	tagURLGenerator gitTagURLGeneratorFunc,
	fileURLGenerator fileURLGeneratorFunc,
) ([]svcmodel.Release, error) {
	r := make([]svcmodel.Release, 0, len(releases))
	for _, release := range releases {
		svcRls, err := ToSvcRelease(release, tagURLGenerator, fileURLGenerator)
		if err != nil {
			return nil, err
		}

		r = append(r, svcRls)
	}

	return r, nil
}

type ReleaseAttachment struct {
	ID        uuid.UUID `json:"attachment_id"`
	Name      string    `json:"name"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}

type fileURLGeneratorFunc func(filePath string) (url.URL, error)

func ToSvcReleaseAttachment(a ReleaseAttachment, fileURLGenerator fileURLGeneratorFunc) (svcmodel.ReleaseAttachment, error) {
	u, err := fileURLGenerator(a.FilePath)
	if err != nil {
		return svcmodel.ReleaseAttachment{}, fmt.Errorf("generating a file URL: %w", err)
	}

	return svcmodel.ReleaseAttachment{
		ID:        a.ID,
		Name:      a.Name,
		FilePath:  a.FilePath,
		URL:       u,
		CreatedAt: a.CreatedAt,
	}, nil
}

func ToSvcReleaseAttachments(
	attachments []ReleaseAttachment,
	fileURLGenerator fileURLGeneratorFunc,
) ([]svcmodel.ReleaseAttachment, error) {
	a := make([]svcmodel.ReleaseAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		svcAttachment, err := ToSvcReleaseAttachment(attachment, fileURLGenerator)
		if err != nil {
			return nil, fmt.Errorf("converting release attachment to service model: %w", err)
		}

		a = append(a, svcAttachment)
	}

	return a, nil
}
