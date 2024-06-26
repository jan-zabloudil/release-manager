package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateEnvironmentInput struct {
	Name       string `json:"name"`
	ServiceURL string `json:"service_url"`
}

type UpdateEnvironmentInput struct {
	Name       *string `json:"name"`
	ServiceURL *string `json:"service_url"`
}

type Environment struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ServiceURL string    `json:"service_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToSvcCreateEnvironmentInput(c CreateEnvironmentInput, projectID uuid.UUID) svcmodel.CreateEnvironmentInput {
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
