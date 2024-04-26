package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateProjectInput struct {
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
}

type UpdateProjectInput struct {
	Name                      *string                              `json:"name"`
	SlackChannelID            *string                              `json:"slack_channel_id"`
	ReleaseNotificationConfig UpdateReleaseNotificationConfigInput `json:"release_notification_config"`
}

type Project struct {
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

type UpdateReleaseNotificationConfigInput struct {
	Message         *string `json:"message"`
	ShowProjectName *bool   `json:"show_project_name"`
	ShowReleaseName *bool   `json:"show_release_name"`
	ShowChangelog   *bool   `json:"show_changelog"`
	ShowDeployments *bool   `json:"show_deployments"`
	ShowSourceCode  *bool   `json:"show_source_code"`
}

func ToSvcCreateProjectInput(c CreateProjectInput) svcmodel.CreateProjectInput {
	return svcmodel.CreateProjectInput{
		Name:                      c.Name,
		SlackChannelID:            c.SlackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(c.ReleaseNotificationConfig),
	}
}

func ToSvcUpdateProjectInput(u UpdateProjectInput) svcmodel.UpdateProjectInput {
	return svcmodel.UpdateProjectInput{
		Name:                            u.Name,
		SlackChannelID:                  u.SlackChannelID,
		ReleaseNotificationConfigUpdate: svcmodel.UpdateReleaseNotificationConfigInput(u.ReleaseNotificationConfig),
	}
}

func ToProject(p svcmodel.Project) Project {
	return Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		CreatedAt:                 p.CreatedAt.Local(),
		UpdatedAt:                 p.UpdatedAt.Local(),
	}
}

func ToProjects(projects []svcmodel.Project) []Project {
	p := make([]Project, 0, len(projects))
	for _, project := range projects {
		p = append(p, ToProject(project))
	}

	return p
}
