package model

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Project struct {
	ID                        uuid.UUID                 `db:"id"`
	Name                      string                    `db:"name"`
	SlackChannelID            string                    `db:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `db:"release_notification_config"`
	GithubRepositoryURL       string                    `db:"github_repository_url"` // TODO remove this field
	GithubOwnerSlug           sql.NullString            `db:"github_owner_slug"`
	GithubRepoSlug            sql.NullString            `db:"github_repo_slug"`
	GithubRepoURL             sql.NullString            `db:"github_repo_url"`
	CreatedAt                 time.Time                 `db:"created_at"`
	UpdatedAt                 time.Time                 `db:"updated_at"`
}

type ReleaseNotificationConfig struct {
	Message          string `json:"message"`
	ShowProjectName  bool   `json:"show_project_name"`
	ShowReleaseTitle bool   `json:"show_release_title"`
	ShowReleaseNotes bool   `json:"show_release_notes"`
	ShowDeployments  bool   `json:"show_deployments"`
	ShowSourceCode   bool   `json:"show_source_code"`
}

func ToSvcProject(p Project) (svcmodel.Project, error) {
	u, err := url.Parse(p.GithubRepositoryURL)
	if err != nil {
		return svcmodel.Project{}, err
	}

	var githubRepo *svcmodel.GithubRepo
	if p.GithubOwnerSlug.Valid && p.GithubRepoSlug.Valid && p.GithubRepoURL.Valid {
		u, err := url.Parse(p.GithubRepoURL.String)
		if err != nil {
			return svcmodel.Project{}, fmt.Errorf("failed to parse github repository URL: %w", err)
		}

		githubRepo = &svcmodel.GithubRepo{
			URL:       *u,
			OwnerSlug: p.GithubOwnerSlug.String,
			RepoSlug:  p.GithubRepoSlug.String,
		}
	}

	return svcmodel.Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepositoryURL:       *u, // TODO remove
		GithubRepo:                githubRepo,
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}, nil
}

func ToSvcProjects(projects []Project) ([]svcmodel.Project, error) {
	p := make([]svcmodel.Project, 0, len(projects))
	for _, project := range projects {
		svcProject, err := ToSvcProject(project)
		if err != nil {
			return nil, err
		}

		p = append(p, svcProject)
	}

	return p, nil
}
