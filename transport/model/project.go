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

type SetProjectGithubRepoInput struct {
	RawRepoURL string `json:"github_repo_url"`
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
	Message          string `json:"message"`
	ShowProjectName  bool   `json:"show_project_name"`
	ShowReleaseTitle bool   `json:"show_release_title"`
	ShowReleaseNotes bool   `json:"show_release_notes"`
	ShowDeployments  bool   `json:"show_deployments"`
	ShowSourceCode   bool   `json:"show_source_code"`
}

type UpdateReleaseNotificationConfigInput struct {
	Message          *string `json:"message"`
	ShowProjectName  *bool   `json:"show_project_name"`
	ShowReleaseTitle *bool   `json:"show_release_title"`
	ShowReleaseNotes *bool   `json:"show_release_notes"`
	ShowDeployments  *bool   `json:"show_deployments"`
	ShowSourceCode   *bool   `json:"show_source_code"`
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
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}

func ToProjects(projects []svcmodel.Project) []Project {
	p := make([]Project, 0, len(projects))
	for _, project := range projects {
		p = append(p, ToProject(project))
	}

	return p
}
