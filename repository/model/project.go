package model

import (
	"net/url"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Project struct {
	ID                        uuid.UUID                 `json:"id"`
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
	GithubRepositoryURL       string                    `json:"github_repository_url"`
	CreatedAt                 time.Time                 `json:"created_at"`
	UpdatedAt                 time.Time                 `json:"updated_at"`
}

// CreateProjectInput is the input used for creating a project and adding an owner as a project member
type CreateProjectInput struct {
	ID                        uuid.UUID                 `json:"p_id"`
	Name                      string                    `json:"p_name"`
	SlackChannelID            string                    `json:"p_slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"p_release_notification_config"`
	GithubRepositoryURL       string                    `json:"p_github_repository"`
	OwnerUserID               uuid.UUID                 `json:"p_user_id"`
	OwnerProjectRole          string                    `json:"p_project_role"`
	CreatedAt                 time.Time                 `json:"p_created_at"`
	UpdatedAt                 time.Time                 `json:"p_updated_at"`
}

type UpdateProjectInput struct {
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
	GithubRepositoryURL       string                    `json:"github_repository_url"`
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

type GithubRepository struct {
	URL            string `json:"url"`
	OwnerSlug      string `json:"owner_slug"`
	RepositorySlug string `json:"repository_slug"`
}

func ToCreateProjectInput(p svcmodel.Project, owner svcmodel.ProjectMember) CreateProjectInput {
	return CreateProjectInput{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepositoryURL:       p.GithubRepositoryURL.String(),
		OwnerUserID:               owner.User.ID,
		OwnerProjectRole:          string(owner.ProjectRole),
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}

func ToUpdateProjectInput(p svcmodel.Project) UpdateProjectInput {
	return UpdateProjectInput{
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepositoryURL:       p.GithubRepositoryURL.String(),
		UpdatedAt:                 p.UpdatedAt,
	}
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
