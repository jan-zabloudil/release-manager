package model

import (
	"net/url"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateEnvironmentRequest struct {
	Name       string `json:"name"`
	ServiceURL string `json:"service_url"`
}

type UpdateEnvironmentRequest struct {
	Name       *string `json:"name"`
	ServiceURL *string `json:"service_url"`
}

type EnvironmentResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ServiceURL string    `json:"service_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToSvcEnvironmentCreation(projectID uuid.UUID, name, rawURL string) svcmodel.EnvironmentCreation {
	return svcmodel.EnvironmentCreation{
		ProjectID:     projectID,
		Name:          name,
		ServiceRawURL: rawURL,
	}
}

func ToSvcEnvironmentUpdate(name, rawURL *string) svcmodel.EnvironmentUpdate {
	return svcmodel.EnvironmentUpdate{
		Name:          name,
		ServiceRawURL: rawURL,
	}
}

func ToEnvironmentResponse(id uuid.UUID, name string, u url.URL, createdAt, updatedAt time.Time) EnvironmentResponse {
	return EnvironmentResponse{
		ID:         id,
		Name:       name,
		ServiceURL: u.String(),
		CreatedAt:  createdAt.Local(),
		UpdatedAt:  updatedAt.Local(),
	}
}

func ToEnvironmentsResponse(envs []svcmodel.Environment) []EnvironmentResponse {
	e := make([]EnvironmentResponse, 0, len(envs))
	for _, env := range envs {
		e = append(e, ToEnvironmentResponse(
			env.ID,
			env.Name,
			env.ServiceURL,
			env.CreatedAt,
			env.UpdatedAt,
		))
	}

	return e
}
