package repository

import (
	"context"

	"release-manager/pkg/dberrors"
	"release-manager/repository/model"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type EnvironmentRepository struct {
	client *supabase.Client
	entity string
}

func NewEnvironmentRepository(c *supabase.Client) *EnvironmentRepository {
	return &EnvironmentRepository{
		client: c,
		entity: "environments",
	}
}

func (r *EnvironmentRepository) Create(ctx context.Context, e svcmodel.Environment) error {
	err := r.client.
		DB.From(r.entity).
		Insert(model.ToEnvironment(
			e.ID,
			e.ProjectID,
			e.Name,
			e.ServiceURL,
			e.CreatedAt,
			e.UpdatedAt,
		)).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *EnvironmentRepository) Read(ctx context.Context, envID uuid.UUID) (svcmodel.Environment, error) {
	var resp model.Environment
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("id", envID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Environment{}, util.ToDBError(err)
	}

	env, err := svcmodel.ToEnvironment(
		resp.ID,
		resp.ProjectID,
		resp.Name,
		resp.ServiceURL,
		resp.CreatedAt,
		resp.UpdatedAt,
	)
	if err != nil {
		return svcmodel.Environment{}, dberrors.NewToSvcModelError().Wrap(err)
	}

	return env, nil
}

func (r *EnvironmentRepository) ReadByNameForProject(ctx context.Context, projectID uuid.UUID, name string) (svcmodel.Environment, error) {
	var resp model.Environment
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("name", name).
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Environment{}, util.ToDBError(err)
	}

	env, err := svcmodel.ToEnvironment(
		resp.ID,
		resp.ProjectID,
		resp.Name,
		resp.ServiceURL,
		resp.CreatedAt,
		resp.UpdatedAt,
	)
	if err != nil {
		return svcmodel.Environment{}, dberrors.NewToSvcModelError().Wrap(err)
	}

	return env, nil
}

func (r *EnvironmentRepository) ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Environment, error) {
	var resp []model.Environment
	err := r.client.
		DB.From(r.entity).
		Select("*").
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	envs, err := model.ToSvcEnvironments(resp)
	if err != nil {
		return nil, dberrors.NewToSvcModelError().Wrap(err)
	}

	return envs, nil
}

func (r *EnvironmentRepository) Delete(ctx context.Context, envID uuid.UUID) error {
	err := r.client.
		DB.From(r.entity).
		Delete().
		Eq("id", envID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *EnvironmentRepository) Update(ctx context.Context, e svcmodel.Environment) error {
	err := r.client.
		DB.From(r.entity).
		Update(model.ToEnvironmentUpdate(
			e.Name,
			e.ServiceURL,
			e.UpdatedAt,
		)).
		Eq("id", e.ID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}
