package model

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"release-manager/pkg/validator"

	"github.com/google/uuid"
)

const (
	// expectedGithubRepositoryURLSlugCount is the expected number of slugs in a GitHub repository URL
	// Example URL: https://github.com/owner/repo -> owner and repo are the slugs
	expectedGithubRepositoryURLSlugCount = 2

	githubHost    = "github.com"
	githubHostWWW = "www.github.com"
)

var (
	errProjectNameRequired               = errors.New("project name is required")
	errProjectInvalidGithubRepositoryURL = errors.New("invalid GitHub repository URL")
)

type Project struct {
	ID                        uuid.UUID
	Name                      string
	SlackChannelID            string
	ReleaseNotificationConfig ReleaseNotificationConfig
	GithubRepository          GithubRepository
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type CreateProjectInput struct {
	Name                      string
	SlackChannelID            string
	ReleaseNotificationConfig ReleaseNotificationConfig
	GithubRepositoryRawURL    string
}

type UpdateProjectInput struct {
	Name                            *string
	SlackChannelID                  *string
	ReleaseNotificationConfigUpdate UpdateReleaseNotificationConfigInput
	GithubRepositoryRawURL          *string
}

type ReleaseNotificationConfig struct {
	Message         string
	ShowProjectName bool
	ShowReleaseName bool
	ShowChangelog   bool
	ShowDeployments bool
	ShowSourceCode  bool
}

type UpdateReleaseNotificationConfigInput struct {
	Message         *string
	ShowProjectName *bool
	ShowReleaseName *bool
	ShowChangelog   *bool
	ShowDeployments *bool
	ShowSourceCode  *bool
}

// GithubRepository represents a GitHub repository URL
// Example URL: https://github.com/owner/repo, OwnerSlug: owner, RepositorySlug: repo
// OwnerSlug and RepositorySlug are saved separately to avoid parsing the URL every time
// Both slugs are needed for the GitHub API
type GithubRepository struct {
	URL            url.URL
	OwnerSlug      string
	RepositorySlug string
}

func NewProject(c CreateProjectInput) (Project, error) {
	repo, err := toGithubRepository(c.GithubRepositoryRawURL)
	if err != nil {
		return Project{}, err
	}

	now := time.Now()
	p := Project{
		ID:                        uuid.New(),
		Name:                      c.Name,
		SlackChannelID:            c.SlackChannelID,
		ReleaseNotificationConfig: c.ReleaseNotificationConfig,
		GithubRepository:          repo,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	if err := p.Validate(); err != nil {
		return Project{}, err
	}

	return p, nil
}

func (p *Project) Update(u UpdateProjectInput) error {
	if u.GithubRepositoryRawURL != nil {
		repo, err := toGithubRepository(*u.GithubRepositoryRawURL)
		if err != nil {
			return err
		}

		p.GithubRepository = repo
	}
	if u.Name != nil {
		p.Name = *u.Name
	}
	if u.SlackChannelID != nil {
		p.SlackChannelID = *u.SlackChannelID
	}

	p.ReleaseNotificationConfig.Update(u.ReleaseNotificationConfigUpdate)
	p.UpdatedAt = time.Now()

	return p.Validate()
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return errProjectNameRequired
	}

	return nil
}

func (p *Project) IsSlackConfigured() bool {
	return p.SlackChannelID != ""
}

func (c *ReleaseNotificationConfig) Update(u UpdateReleaseNotificationConfigInput) {
	if u.Message != nil {
		c.Message = *u.Message
	}
	if u.ShowProjectName != nil {
		c.ShowProjectName = *u.ShowProjectName
	}
	if u.ShowReleaseName != nil {
		c.ShowReleaseName = *u.ShowReleaseName
	}
	if u.ShowChangelog != nil {
		c.ShowChangelog = *u.ShowChangelog
	}
	if u.ShowDeployments != nil {
		c.ShowDeployments = *u.ShowDeployments
	}
	if u.ShowSourceCode != nil {
		c.ShowSourceCode = *u.ShowSourceCode
	}
}

// toGithubRepository parses a GitHub repository URL and returns a GithubRepository struct
// Github Repository URL must be in the format: http(s)://github.com/owner/repo or http(s)://www.github.com/owner/repo
func toGithubRepository(rawURL string) (GithubRepository, error) {
	if rawURL == "" {
		return GithubRepository{}, nil
	}

	if !validator.IsAbsoluteURL(rawURL) {
		return GithubRepository{}, errProjectInvalidGithubRepositoryURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return GithubRepository{}, errProjectInvalidGithubRepositoryURL
	}

	if u.Host != githubHost && u.Host != githubHostWWW {
		return GithubRepository{}, errProjectInvalidGithubRepositoryURL
	}

	path := strings.Trim(u.Path, "/")
	slugs := strings.Split(path, "/")

	if len(slugs) != expectedGithubRepositoryURLSlugCount {
		return GithubRepository{}, errProjectInvalidGithubRepositoryURL
	}

	if slugs[0] == "" || slugs[1] == "" {
		return GithubRepository{}, errProjectInvalidGithubRepositoryURL
	}

	return GithubRepository{
		URL:            *u,
		OwnerSlug:      slugs[0],
		RepositorySlug: slugs[1],
	}, nil
}
