package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type App struct {
	ID           uuid.UUID    `json:"id"`
	ProjectID    uuid.UUID    `json:"project_id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Environments Environments `json:"environments"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type Environments struct {
	DevURL string `json:"dev_url,omitempty"`
	StgURL string `json:"stg_url,omitempty"`
	PrdURL string `json:"prd_url,omitempty"`
}

func ToDDApp(id, projectID uuid.UUID, name, description string, env svcmodel.Environments, createdAt, updatedAt time.Time) App {
	var dbEnv Environments

	if env.DevURL != nil {
		dbEnv.DevURL = env.DevURL.String()
	}
	if env.StgURL != nil {
		dbEnv.StgURL = env.StgURL.String()
	}
	if env.PrdURL != nil {
		dbEnv.PrdURL = env.PrdURL.String()
	}

	return App{
		ID:           id,
		ProjectID:    projectID,
		Name:         name,
		Description:  description,
		Environments: dbEnv,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func ToSvcApp(id, projectID uuid.UUID, name, description string, env Environments, createdAt, updatedAt time.Time) (svcmodel.App, error) {
	var svcEnv svcmodel.Environments
	var err error

	svcEnv.DevURL, err = toSvcEnvURL(env.DevURL)
	if err != nil {
		return svcmodel.App{}, err
	}

	svcEnv.StgURL, err = toSvcEnvURL(env.StgURL)
	if err != nil {
		return svcmodel.App{}, err
	}

	svcEnv.PrdURL, err = toSvcEnvURL(env.PrdURL)
	if err != nil {
		return svcmodel.App{}, err
	}

	return svcmodel.App{
		ID:           id,
		ProjectID:    projectID,
		Name:         name,
		Description:  description,
		Environments: svcEnv,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func ToSvcApps(apps []App) ([]svcmodel.App, error) {
	a := make([]svcmodel.App, 0, len(apps))
	for _, dbApp := range apps {
		app, err := ToSvcApp(
			dbApp.ID,
			dbApp.ProjectID,
			dbApp.Name,
			dbApp.Description,
			dbApp.Environments,
			dbApp.CreatedAt,
			dbApp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		a = append(a, app)
	}

	return a, nil
}

func toSvcEnvURL(rawURL string) (svcmodel.EnvURL, error) {
	if rawURL == "" {
		return nil, nil
	}
	return svcmodel.NewEnvURL(rawURL)
}
