package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	errProjectNameRequired = errors.New("project name is required")
)

type Project struct {
	ID                        uuid.UUID
	Name                      string
	SlackChannelID            string
	ReleaseNotificationConfig ReleaseNotificationConfig
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type ProjectCreation struct {
	Name                      string
	SlackChannelID            string
	ReleaseNotificationConfig ReleaseNotificationConfig
}

type ProjectUpdate struct {
	Name                            *string
	SlackChannelID                  *string
	ReleaseNotificationConfigUpdate ReleaseNotificationConfigUpdate
}

type ReleaseNotificationConfig struct {
	Message         string
	ShowProjectName bool
	ShowReleaseName bool
	ShowChangelog   bool
	ShowDeployments bool
	ShowSourceCode  bool
}

type ReleaseNotificationConfigUpdate struct {
	Message         *string
	ShowProjectName *bool
	ShowReleaseName *bool
	ShowChangelog   *bool
	ShowDeployments *bool
	ShowSourceCode  *bool
}

func NewProject(c ProjectCreation) (Project, error) {
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

func (p *Project) Update(u ProjectUpdate) error {
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

func ToProject(
	id uuid.UUID,
	name,
	slackChannelID string,
	rlsCfg ReleaseNotificationConfig,
	createdAt,
	updatedAt time.Time,
) (Project, error) {
	p := Project{
		ID:                        id,
		Name:                      name,
		SlackChannelID:            slackChannelID,
		ReleaseNotificationConfig: rlsCfg,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}

	if err := p.Validate(); err != nil {
		return Project{}, err
	}

	return p, nil
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

func (c *ReleaseNotificationConfig) Update(u ReleaseNotificationConfigUpdate) {
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
