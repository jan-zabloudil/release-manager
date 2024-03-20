package model

import (
	"context"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type AppService interface {
	Create(ctx context.Context, app svcmodel.App) (svcmodel.App, error)
	Get(ctx context.Context, id uuid.UUID) (svcmodel.App, error)
	GetAllForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.App, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, app svcmodel.App) (svcmodel.App, error)
}

type App struct {
	ID           uuid.UUID     `json:"id"`
	Name         *string       `json:"name" validate:"required"`
	Description  *string       `json:"description"`
	Environments *Environments `json:"environments"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type AppPatch struct {
	Name         *string       `json:"name"`
	Description  *string       `json:"description"`
	Environments *Environments `json:"environments"`
}

type Environments struct {
	DevURL *string `json:"dev_url,omitempty" validate:"omitempty,env_url"`
	StgURL *string `json:"stg_url,omitempty" validate:"omitempty,env_url"`
	PrdURL *string `json:"prd_url,omitempty" validate:"omitempty,env_url"`
}

func ToSvcApp(app svcmodel.App, name, description *string, env *Environments) (svcmodel.App, error) {

	if name != nil {
		app.Name = *name
	}
	if description != nil {
		app.Description = *description
	}

	if env != nil {
		netUrls := []*string{env.DevURL, env.StgURL, env.PrdURL}
		envURLs := []*svcmodel.EnvURL{&app.Environments.DevURL, &app.Environments.StgURL, &app.Environments.PrdURL}

		for i, netUrl := range netUrls {
			if netUrl != nil {
				envURL, err := svcmodel.NewEnvURL(*netUrl)
				if err != nil {
					return svcmodel.App{}, err
				}
				*envURLs[i] = envURL
			}
		}
	}
	app.UpdatedAt = time.Now()

	return app, nil
}

func NewSvcApp(projectID uuid.UUID, name, description *string, env *Environments) (svcmodel.App, error) {
	app := svcmodel.App{
		ID:        uuid.New(),
		ProjectID: projectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return ToSvcApp(app, name, description, env)
}

func ToNetApp(id uuid.UUID, name, description string, env svcmodel.Environments, createdAt, updatedAt time.Time) App {
	return App{
		ID:          id,
		Name:        &name,
		Description: &description,
		Environments: &Environments{
			DevURL: toNetEnvURL(env.DevURL),
			StgURL: toNetEnvURL(env.StgURL),
			PrdURL: toNetEnvURL(env.PrdURL),
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func ToNetApps(svcApps []svcmodel.App) []App {
	a := make([]App, 0, len(svcApps))
	for _, svcApp := range svcApps {
		a = append(a, ToNetApp(
			svcApp.ID,
			svcApp.Name,
			svcApp.Description,
			svcApp.Environments,
			svcApp.CreatedAt,
			svcApp.UpdatedAt,
		))
	}

	return a
}

func toNetEnvURL(url svcmodel.EnvURL) *string {
	if url != nil {
		urlString := url.String()
		return &urlString
	}
	return nil
}
