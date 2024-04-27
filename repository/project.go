package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type ProjectRepository struct {
	client *supabase.Client
	entity string
}

func NewProjectRepository(c *supabase.Client) *ProjectRepository {
	return &ProjectRepository{
		client: c,
		entity: "projects",
	}
}

func (r *ProjectRepository) Create(ctx context.Context, p svcmodel.Project) error {
	data := model.ToProject(p)

	err := r.client.
		DB.From(r.entity).
		Insert(&data).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.Project, error) {
	var resp model.Project
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Project{}, util.ToDBError(err)
	}

	return model.ToSvcProject(resp), nil
}

func (r *ProjectRepository) ReadAll(ctx context.Context) ([]svcmodel.Project, error) {
	var resp []model.Project
	err := r.client.
		DB.From(r.entity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	return model.ToSvcProjects(resp), nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.
		DB.From(r.entity).
		Delete().Eq("id", id.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) Update(ctx context.Context, p svcmodel.Project) error {
	data := model.ToProjectUpdate(p)

	err := r.client.
		DB.From(r.entity).
		Update(&data).
		Eq("id", (p.ID).String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}
