package model

import (
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
	GithubRepositoryURL       string                    `db:"github_repository_url"`
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

	return svcmodel.Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepositoryURL:       *u,
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
