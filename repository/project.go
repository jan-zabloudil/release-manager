package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
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

func (r *ProjectRepository) Insert(ctx context.Context, p svcmodel.Project) (svcmodel.Project, error) {

	var resp []model.ProjectResponse
	err := r.client.
		DB.From(r.entity).
		Insert(model.ToProjectDBInput(p)).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Project{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetch(resp); err != nil {
		return svcmodel.Project{}, err
	}

	return model.ToSvcProject(resp[0]), nil
}

func (r *ProjectRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.Project, error) {

	var resp model.ProjectResponse
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Project{}, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcProject(resp), nil
}

func (r *ProjectRepository) ReadAll(ctx context.Context) ([]svcmodel.Project, error) {

	var resp []model.ProjectResponse
	err := r.client.
		DB.From(r.entity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcProjects(resp), nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.
		DB.From(r.entity).
		Delete().Eq("id", id.String()).ExecuteWithContext(ctx, nil)
	if err != nil {
		return utils.WrapSupabaseDBErr(err)
	}

	return nil
}

func (r *ProjectRepository) Update(ctx context.Context, p svcmodel.Project) (svcmodel.Project, error) {

	var resp []model.ProjectResponse
	err := r.client.
		DB.From(r.entity).
		Update(model.ToProjectDBInput(p)).
		Eq("id", (p.ID).String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Project{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetch(resp); err != nil {
		return svcmodel.Project{}, err
	}

	return model.ToSvcProject(resp[0]), nil
}
