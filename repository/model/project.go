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

func ToProject(
	id uuid.UUID,
	name,
	slackChannelID string,
	rlsCfg svcmodel.ReleaseNotificationConfig,
	createdAt,
	updatedAt time.Time,
) Project {
	return Project{
		ID:                        id,
		Name:                      name,
		SlackChannelID:            slackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(rlsCfg),
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}
}

func ToProjectUpdate(
	name,
	slackChannelID string,
	rlsCfg svcmodel.ReleaseNotificationConfig,
	updatedAt time.Time,
) ProjectUpdate {
	return ProjectUpdate{
		Name:                      name,
		SlackChannelID:            slackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(rlsCfg),
		UpdatedAt:                 updatedAt,
	}
}

func ToSvcProject(
	id uuid.UUID,
	name,
	slackChannelID string,
	cfg ReleaseNotificationConfig,
	createdAt,
	updatedAt time.Time,
) (svcmodel.Project, error) {
	return svcmodel.ToProject(
		id,
		name,
		slackChannelID,
		svcmodel.ReleaseNotificationConfig(cfg),
		createdAt,
		updatedAt,
	)
}

func ToSvcProjects(projects []Project) ([]svcmodel.Project, error) {
	svcProjects := make([]svcmodel.Project, 0, len(projects))
	for _, p := range projects {
		svcProject, err := ToSvcProject(
			p.ID,
			p.Name,
			p.SlackChannelID,
			p.ReleaseNotificationConfig,
			p.CreatedAt,
			p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		svcProjects = append(svcProjects, svcProject)
	}

	return svcProjects, nil
}
