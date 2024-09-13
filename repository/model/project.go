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
	GithubOwnerSlug           sql.NullString            `db:"github_owner_slug"`
	GithubRepoSlug            sql.NullString            `db:"github_repo_slug"`
	CreatedAt                 time.Time                 `db:"created_at"`
	UpdatedAt                 time.Time                 `db:"updated_at"`
}

type ReleaseNotificationConfig struct {
	Message            string `json:"message"`
	ShowProjectName    bool   `json:"show_project_name"`
	ShowReleaseTitle   bool   `json:"show_release_title"`
	ShowReleaseNotes   bool   `json:"show_release_notes"`
	ShowLastDeployment bool   `json:"show_last_deployment"`
	ShowSourceCode     bool   `json:"show_source_code"`
}

type githubRepoURLGeneratorFunc func(ownerSlug, repoSlug string) (url.URL, error)

func ToSvcProject(p Project, urlGenerator githubRepoURLGeneratorFunc) (svcmodel.Project, error) {
	var repo *svcmodel.GithubRepo
	if p.GithubOwnerSlug.Valid && p.GithubRepoSlug.Valid {
		repoURL, err := urlGenerator(p.GithubOwnerSlug.String, p.GithubRepoSlug.String)
		if err != nil {
			return svcmodel.Project{}, fmt.Errorf("generating repo URL: %w", err)
		}

		repo = &svcmodel.GithubRepo{
			OwnerSlug: p.GithubOwnerSlug.String,
			RepoSlug:  p.GithubRepoSlug.String,
			URL:       repoURL,
		}
	}

	return svcmodel.Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepo:                repo,
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}, nil
}

func ToSvcProjects(projects []Project, urlGenerator githubRepoURLGeneratorFunc) ([]svcmodel.Project, error) {
	p := make([]svcmodel.Project, 0, len(projects))
	for _, project := range projects {
		svcProject, err := ToSvcProject(project, urlGenerator)
		if err != nil {
			return nil, err
		}

		p = append(p, svcProject)
	}

	return p, nil
}
