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
	errReleaseNotificationConfigMessageRequired = errors.New("message in release notification config is required")
)

type Project struct {
	ID                        uuid.UUID
	Name                      string
	SlackChannelID            string
	ReleaseNotificationConfig ReleaseNotificationConfig
	GithubRepo                *GithubRepo
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type GithubRepo struct {
	URL       url.URL
	OwnerSlug string
	RepoSlug  string
}

type CreateProjectInput struct {
	Name                      string
	SlackChannelID            string
	ReleaseNotificationConfig ReleaseNotificationConfig
}

type UpdateProjectInput struct {
	Name                            *string
	SlackChannelID                  *string
	ReleaseNotificationConfigUpdate UpdateReleaseNotificationConfigInput
}

type ReleaseNotificationConfig struct {
	Message            string
	ShowProjectName    bool
	ShowReleaseTitle   bool
	ShowReleaseNotes   bool
	ShowLastDeployment bool
	ShowSourceCode     bool
}

type UpdateReleaseNotificationConfigInput struct {
	Message            *string
	ShowProjectName    *bool
	ShowReleaseTitle   *bool
	ShowReleaseNotes   *bool
	ShowLastDeployment *bool
	ShowSourceCode     *bool
}

func NewProject(c CreateProjectInput) (Project, error) {
	now := time.Now()
	p := Project{
		ID:                        uuid.New(),
		Name:                      c.Name,
		SlackChannelID:            c.SlackChannelID,
		ReleaseNotificationConfig: c.ReleaseNotificationConfig,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	if err := p.Validate(); err != nil {
		return Project{}, err
	}

	return p, nil
}

func (p *Project) SetGithubRepo(repo *GithubRepo) {
	p.GithubRepo = repo
	p.UpdatedAt = time.Now()
}

func (p *Project) Update(u UpdateProjectInput) error {
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

func (p *Project) IsGithubRepoSet() bool {
	return p.GithubRepo != nil
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
	if u.ShowLastDeployment != nil {
		c.ShowLastDeployment = *u.ShowLastDeployment
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
