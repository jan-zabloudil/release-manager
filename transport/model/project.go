package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
}

type UpdateProjectRequest struct {
	Name                      *string                                 `json:"name"`
	SlackChannelID            *string                                 `json:"slack_channel_id"`
	ReleaseNotificationConfig *UpdateReleaseNotificationConfigRequest `json:"release_notification_config"`
}

type ProjectResponse struct {
	ID                        uuid.UUID                 `json:"id"`
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
	CreatedAt                 time.Time                 `json:"created_at"`
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

type UpdateReleaseNotificationConfigRequest struct {
	Message         *string `json:"message"`
	ShowProjectName *bool   `json:"show_project_name"`
	ShowReleaseName *bool   `json:"show_release_name"`
	ShowChangelog   *bool   `json:"show_changelog"`
	ShowDeployments *bool   `json:"show_deployments"`
	ShowSourceCode  *bool   `json:"show_source_code"`
}

func ToSvcProjectCreation(name, slackChannelID string, rlsCfg ReleaseNotificationConfig) svcmodel.ProjectCreation {
	return svcmodel.ProjectCreation{
		Name:                      name,
		SlackChannelID:            slackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(rlsCfg),
	}
}

func ToSvcProjectUpdate(
	name,
	slackChannelID *string,
	rlsCfg *UpdateReleaseNotificationConfigRequest,
) svcmodel.ProjectUpdate {
	var p svcmodel.ProjectUpdate

	p.Name, p.SlackChannelID = name, slackChannelID
	if rlsCfg != nil {
		p.ReleaseNotificationConfigUpdate = svcmodel.ReleaseNotificationConfigUpdate(*rlsCfg)
	}

	return p
}

func ToProjectResponse(
	id uuid.UUID,
	name,
	slackChannelID string,
	rlsConfig svcmodel.ReleaseNotificationConfig,
	createdAt,
	updatedAt time.Time,
) ProjectResponse {
	return ProjectResponse{
		ID:                        id,
		Name:                      name,
		SlackChannelID:            slackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(rlsConfig),
		CreatedAt:                 createdAt.Local(),
		UpdatedAt:                 updatedAt.Local(),
	}
}

func ToProjects(projects []svcmodel.Project) []ProjectResponse {
	p := make([]ProjectResponse, 0, len(projects))
	for _, project := range projects {
		p = append(p, ToProjectResponse(
			project.ID,
			project.Name,
			project.SlackChannelID,
			project.ReleaseNotificationConfig,
			project.CreatedAt,
			project.UpdatedAt,
		))
	}

	return p
}
