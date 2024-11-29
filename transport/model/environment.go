package model

import (
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
)

type CreateEnvironmentInput struct {
	Name       string `json:"name" validate:"required"`
	ServiceURL string `json:"service_url" validate:"omitempty,http_url"`
}

type UpdateEnvironmentInput struct {
	Name       *string `json:"name" validate:"omitempty,min=1"`
	ServiceURL *string `json:"service_url" validate:"omitempty,optional_http_url"`
}

type Environment struct {
	ID         id.Environment `json:"id"`
	Name       string         `json:"name"`
	ServiceURL string         `json:"service_url"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type EnvironmentURLParams struct {
	ProjectID     id.Project     `param:"path=project_id"`
	EnvironmentID id.Environment `param:"path=environment_id"`
}

func ToSvcCreateEnvironmentInput(c CreateEnvironmentInput, projectID id.Project) svcmodel.CreateEnvironmentInput {
	return svcmodel.CreateEnvironmentInput{
		ProjectID:     projectID,
		Name:          c.Name,
		ServiceRawURL: c.ServiceURL,
	}
}

func ToSvcUpdateEnvironmentInput(u UpdateEnvironmentInput) svcmodel.UpdateEnvironmentInput {
	return svcmodel.UpdateEnvironmentInput{
		Name:          u.Name,
		ServiceRawURL: u.ServiceURL,
	}
}

func ToEnvironment(e svcmodel.Environment) Environment {
	return Environment{
		ID:         e.ID,
		Name:       e.Name,
		ServiceURL: e.ServiceURL.String(),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func ToEnvironments(envs []svcmodel.Environment) []Environment {
	e := make([]Environment, 0, len(envs))
	for _, env := range envs {
		e = append(e, ToEnvironment(env))
	}

	return e
}
