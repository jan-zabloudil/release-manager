package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Project struct {
	ID                        uuid.UUID                 `json:"id"`
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
	CreatedAt                 time.Time                 `json:"created_at"`
	UpdatedAt                 time.Time                 `json:"updated_at"`
}

type ProjectUpdate struct {
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
	UpdatedAt                 time.Time                 `json:"updated_at"`
}

type ReleaseNotificationConfig struct {
	Message         string `json:"message"`
	ShowProjectName bool   `json:"show_project_name"`
	ShowReleaseName bool   `json:"show_release_name"`
	ShowChangelog   bool   `json:"show_changelog"`
	ShowDeployments bool   `json:"show_deployments"`
	ShowSourceCode  bool   `json:"show_source_code"`
}

func ToProject(p svcmodel.Project) Project {
	return Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}

func ToProjectUpdate(p svcmodel.Project) ProjectUpdate {
	return ProjectUpdate{
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		UpdatedAt:                 p.UpdatedAt,
	}
}

func ToSvcProject(p Project) svcmodel.Project {
	return svcmodel.Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}

func ToSvcProjects(projects []Project) []svcmodel.Project {
	p := make([]svcmodel.Project, 0, len(projects))
	for _, project := range projects {
		p = append(p, ToSvcProject(project))
	}

	return p
}
