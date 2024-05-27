package model

import (
	"errors"
	"net/url"
	"time"

	cryptox "release-manager/pkg/crypto"

	"github.com/google/uuid"
)

var (
	errProjectNameRequired                      = errors.New("project name is required")
	errProjectGithubRepoURLCannotBeParsed       = errors.New("github repository URL cannot be parsed")
	errReleaseNotificationConfigMessageRequired = errors.New("message in release notification config is required")
)

type Project struct {
	ID                        uuid.UUID
	Name                      string
	SlackChannelID            string
	ReleaseNotificationConfig ReleaseNotificationConfig
	GithubRepositoryURL       url.URL
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
	Message          string
	ShowProjectName  bool
	ShowReleaseTitle bool
	ShowReleaseNotes bool
	ShowDeployments  bool
	ShowSourceCode   bool
}

type UpdateReleaseNotificationConfigInput struct {
	Message          *string
	ShowProjectName  *bool
	ShowReleaseTitle *bool
	ShowReleaseNotes *bool
	ShowDeployments  *bool
	ShowSourceCode   *bool
}

func NewProject(c CreateProjectInput) (Project, error) {
	u, err := url.Parse(c.GithubRepositoryRawURL)
	if err != nil {
		return Project{}, errProjectGithubRepoURLCannotBeParsed
	}

	now := time.Now()
	p := Project{
		ID:                        uuid.New(),
		Name:                      c.Name,
		SlackChannelID:            c.SlackChannelID,
		ReleaseNotificationConfig: c.ReleaseNotificationConfig,
		GithubRepositoryURL:       *u,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	if err := p.Validate(); err != nil {
		return Project{}, err
	}

	return p, nil
}

type UpdateProjectFunc func(p Project) (Project, error)

func (p *Project) Update(u UpdateProjectInput) error {
	if u.GithubRepositoryRawURL != nil {
		u, err := url.Parse(*u.GithubRepositoryRawURL)
		if err != nil {
			return errProjectGithubRepoURLCannotBeParsed
		}

		p.GithubRepositoryURL = *u
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

	return p.ReleaseNotificationConfig.Validate()
}

func (p *Project) IsSlackChannelSet() bool {
	return p.SlackChannelID != ""
}

func (p *Project) IsGithubConfigured() bool {
	return p.GithubRepositoryURL != (url.URL{})
}

func (c *ReleaseNotificationConfig) Update(u UpdateReleaseNotificationConfigInput) {
	if u.Message != nil {
		c.Message = *u.Message
	}
	if u.ShowProjectName != nil {
		c.ShowProjectName = *u.ShowProjectName
	}
	if u.ShowReleaseTitle != nil {
		c.ShowReleaseTitle = *u.ShowReleaseTitle
	}
	if u.ShowReleaseNotes != nil {
		c.ShowReleaseNotes = *u.ShowReleaseNotes
	}
	if u.ShowDeployments != nil {
		c.ShowDeployments = *u.ShowDeployments
	}
	if u.ShowSourceCode != nil {
		c.ShowSourceCode = *u.ShowSourceCode
	}
}

func (c *ReleaseNotificationConfig) IsEmpty() bool {
	if c == nil {
		return true
	}

	return *c == ReleaseNotificationConfig{}
}

func (c *ReleaseNotificationConfig) Validate() error {
	if c.Message == "" {
		return errReleaseNotificationConfigMessageRequired
	}

	return nil
}

type ProjectInvitationEmailData struct {
	ProjectName string
	Token       string
}

func NewProjectInvitationEmailData(projectName string, token cryptox.Token) ProjectInvitationEmailData {
	return ProjectInvitationEmailData{
		ProjectName: projectName,
		Token:       string(token),
	}
}
