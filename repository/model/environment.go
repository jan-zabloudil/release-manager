package model

import (
	"net/url"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Environment struct {
	ID         uuid.UUID `db:"id"`
	ProjectID  uuid.UUID `db:"project_id"`
	Name       string    `db:"name"`
	ServiceURL string    `db:"service_url"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func ToSvcEnvironment(e Environment) (svcmodel.Environment, error) {
	u, err := url.Parse(e.ServiceURL)
	if err != nil {
		return svcmodel.Environment{}, err
	}

	return svcmodel.Environment{
		ID:         e.ID,
		ProjectID:  e.ProjectID,
		Name:       e.Name,
		ServiceURL: *u,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}, nil
}

func ToSvcEnvironments(envs []Environment) ([]svcmodel.Environment, error) {
	svcEnvs := make([]svcmodel.Environment, 0, len(envs))
	for _, e := range envs {
		svcEnv, err := ToSvcEnvironment(e)
		if err != nil {
			return nil, err
		}

		svcEnvs = append(svcEnvs, svcEnv)
	}

	return svcEnvs, nil
}
