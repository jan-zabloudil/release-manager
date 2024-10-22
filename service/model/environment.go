package model

import (
	"errors"
	"net/url"
	"time"

	"release-manager/pkg/id"
	"release-manager/pkg/validator"
)

var (
	errEnvironmentInvalidServiceURL        = errors.New("invalid service url")
	errEnvironmentServiceURLMustBeAbsolute = errors.New("service url must be absolute")
	errEnvironmentNameRequired             = errors.New("environment name is required")
)

type Environment struct {
	ID         id.Environment
	ProjectID  id.Project
	Name       string
	ServiceURL url.URL
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateEnvironmentInput struct {
	ProjectID     id.Project
	Name          string
	ServiceRawURL string
}

type UpdateEnvironmentInput struct {
	Name          *string
	ServiceRawURL *string
}

func NewEnvironment(c CreateEnvironmentInput) (Environment, error) {
	u, err := toServiceURL(c.ServiceRawURL)
	if err != nil {
		return Environment{}, err
	}

	now := time.Now()
	env := Environment{
		ID:         id.NewEnvironment(),
		ProjectID:  c.ProjectID,
		Name:       c.Name,
		ServiceURL: u,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := env.Validate(); err != nil {
		return Environment{}, err
	}

	return env, nil
}

func (e *Environment) Validate() error {
	if e.Name == "" {
		return errEnvironmentNameRequired
	}

	if e.ServiceURL.String() != "" && !validator.IsAbsoluteURL(e.ServiceURL.String()) {
		return errEnvironmentServiceURLMustBeAbsolute
	}

	return nil
}

func (e *Environment) Update(u UpdateEnvironmentInput) error {
	if u.ServiceRawURL != nil {
		svcURL, err := toServiceURL(*u.ServiceRawURL)
		if err != nil {
			return err
		}

		e.ServiceURL = svcURL
	}

	if u.Name != nil {
		e.Name = *u.Name
	}

	e.UpdatedAt = time.Now()

	return e.Validate()
}

func (e *Environment) IsServiceURLSet() bool {
	return e.ServiceURL.String() != ""
}

func toServiceURL(rawURL string) (url.URL, error) {
	if rawURL == "" {
		return url.URL{}, nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, errEnvironmentInvalidServiceURL
	}

	return *u, nil
}
