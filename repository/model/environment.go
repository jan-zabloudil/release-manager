package model

import (
	"net/url"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Environment struct {
	ID         uuid.UUID `json:"id"`
	ProjectID  uuid.UUID `json:"project_id"`
	Name       string    `json:"name"`
	ServiceURL string    `json:"service_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type EnvironmentUpdate struct {
	Name       string    `json:"name"`
	ServiceURL string    `json:"service_url"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToEnvironment(id, projectID uuid.UUID, name string, u url.URL, createdAt, updatedAt time.Time) Environment {
	return Environment{
		ID:         id,
		ProjectID:  projectID,
		Name:       name,
		ServiceURL: u.String(),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func ToEnvironmentUpdate(name string, serviceURL url.URL, updatedAt time.Time) EnvironmentUpdate {
	return EnvironmentUpdate{
		Name:       name,
		ServiceURL: serviceURL.String(),
		UpdatedAt:  updatedAt,
	}
}

func ToSvcEnvironments(envs []Environment) ([]svcmodel.Environment, error) {
	svcEnvs := make([]svcmodel.Environment, 0, len(envs))
	for _, e := range envs {
		svcEnv, err := svcmodel.ToEnvironment(
			e.ID,
			e.ProjectID,
			e.Name,
			e.ServiceURL,
			e.CreatedAt,
			e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		svcEnvs = append(svcEnvs, svcEnv)
	}

	return svcEnvs, nil
}
