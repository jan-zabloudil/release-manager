package model

import (
	"database/sql"
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

func ToSvcProject(p Project, repo *svcmodel.GithubRepo) svcmodel.Project {
	return svcmodel.Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepo:                repo,
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}
