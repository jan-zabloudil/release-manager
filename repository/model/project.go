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
	GithubRepository          GithubRepository          `json:"github_repository"`
	CreatedAt                 time.Time                 `json:"created_at"`
	UpdatedAt                 time.Time                 `json:"updated_at"`
}

type UpdateProjectInput struct {
	Name                      string                    `json:"name"`
	SlackChannelID            string                    `json:"slack_channel_id"`
	ReleaseNotificationConfig ReleaseNotificationConfig `json:"release_notification_config"`
	GithubRepository          GithubRepository          `json:"github_repository"`
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

type GithubRepository struct {
	URL            string `json:"url"`
	OwnerSlug      string `json:"owner_slug"`
	RepositorySlug string `json:"repository_slug"`
}

func ToProject(p svcmodel.Project) Project {
	return Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepository:          ToGithubRepository(p.GithubRepository),
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}

func ToUpdateProjectInput(p svcmodel.Project) UpdateProjectInput {
	return UpdateProjectInput{
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepository:          ToGithubRepository(p.GithubRepository),
		UpdatedAt:                 p.UpdatedAt,
	}
}

func ToSvcProject(p Project) (svcmodel.Project, error) {
	repo, err := ToSvcGithubRepository(p.GithubRepository)
	if err != nil {
		return svcmodel.Project{}, err
	}

	return svcmodel.Project{
		ID:                        p.ID,
		Name:                      p.Name,
		SlackChannelID:            p.SlackChannelID,
		ReleaseNotificationConfig: svcmodel.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
		GithubRepository:          repo,
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

func ToGithubRepository(r svcmodel.GithubRepository) GithubRepository {
	return GithubRepository{
		URL:            r.URL.String(),
		OwnerSlug:      r.OwnerSlug,
		RepositorySlug: r.RepositorySlug,
	}
}

func ToSvcGithubRepository(r GithubRepository) (svcmodel.GithubRepository, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return svcmodel.GithubRepository{}, err
	}

	return svcmodel.GithubRepository{
		URL:            *u,
		OwnerSlug:      r.OwnerSlug,
		RepositorySlug: r.RepositorySlug,
	}, nil
}
