package model

import (
	"errors"
	"net/url"
	"time"

	urlx "release-manager/pkg/url"

	"github.com/google/uuid"
)

var (
	errEnvironmentInvalidServiceURL        = errors.New("invalid service url")
	errEnvironmentServiceURLMustBeAbsolute = errors.New("service url must be absolute")
	errEnvironmentNameRequired             = errors.New("environment name is required")
)

type Environment struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	Name       string
	ServiceURL url.URL
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// TODO add release id once releases are implemented
}

type EnvironmentCreation struct {
	ProjectID     uuid.UUID
	Name          string
	ServiceRawURL string
}

type EnvironmentUpdate struct {
	Name          *string
	ServiceRawURL *string
}

func NewEnvironment(c EnvironmentCreation) (Environment, error) {
	u, err := toServiceURL(c.ServiceRawURL)
	if err != nil {
		return Environment{}, err
	}

	now := time.Now()
	env := Environment{
		ID:         uuid.New(),
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

	return nil
}

func ToEnvironment(id, projectID uuid.UUID, name, url string, createdAt, updatedAt time.Time) (Environment, error) {
	u, err := toServiceURL(url)
	if err != nil {
		return Environment{}, err
	}

	env := Environment{
		ID:         id,
		ProjectID:  projectID,
		Name:       name,
		ServiceURL: u,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	if err := env.Validate(); err != nil {
		return Environment{}, err
	}

	return env, nil
}

func (e *Environment) Update(u EnvironmentUpdate) error {
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

func toServiceURL(rawURL string) (url.URL, error) {
	if rawURL == "" {
		return url.URL{}, nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, errEnvironmentInvalidServiceURL
	}

	if !urlx.IsAbsolute(rawURL) {
		return url.URL{}, errEnvironmentServiceURLMustBeAbsolute
	}

	return *u, nil
}
