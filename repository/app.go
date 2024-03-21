package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type AppRepository struct {
	client *supabase.Client
	entity string
}

func NewAppRepository(c *supabase.Client) *AppRepository {
	return &AppRepository{
		client: c,
		entity: "apps",
	}
}

func (r *AppRepository) Insert(ctx context.Context, app svcmodel.App) (svcmodel.App, error) {

	var resp []model.App
	err := r.client.
		DB.From(r.entity).
		Insert(model.ToDDApp(
			app.ID,
			app.ProjectID,
			app.Name,
			app.Description,
			app.Environments,
			app.CreatedAt,
			app.UpdatedAt,
		)).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.App{}, err
	}

	if err := utils.ValidateSingleRecordFetchAfterWriteOperation(resp); err != nil {
		return svcmodel.App{}, err
	}

	return model.ToSvcApp(
		resp[0].ID,
		resp[0].ProjectID,
		resp[0].Name,
		resp[0].Description,
		resp[0].Environments,
		resp[0].CreatedAt,
		resp[0].UpdatedAt,
	)
}

func (r *AppRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.App, error) {

	var resp model.App
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.App{}, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcApp(
		resp.ID,
		resp.ProjectID,
		resp.Name,
		resp.Description,
		resp.Environments,
		resp.CreatedAt,
		resp.UpdatedAt,
	)
}

func (r *AppRepository) ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.App, error) {

	var resp []model.App
	err := r.client.
		DB.From(r.entity).
		Select("*").
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcApps(resp)
}

func (r *AppRepository) Delete(ctx context.Context, id uuid.UUID) error {

	err := r.client.
		DB.From(r.entity).
		Delete().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return utils.WrapSupabaseDBErr(err)
	}

	return nil
}

func (r *AppRepository) Update(ctx context.Context, app svcmodel.App) (svcmodel.App, error) {

	var resp []model.App
	err := r.client.
		DB.From(r.entity).
		Update(model.ToDDApp(
			app.ID,
			app.ProjectID,
			app.Name,
			app.Description,
			app.Environments,
			app.CreatedAt,
			app.UpdatedAt,
		)).
		Eq("id", app.ID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.App{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetchAfterWriteOperation(resp); err != nil {
		return svcmodel.App{}, err
	}

	return model.ToSvcApp(
		resp[0].ID,
		resp[0].ProjectID,
		resp[0].Name,
		resp[0].Description,
		resp[0].Environments,
		resp[0].CreatedAt,
		resp[0].UpdatedAt,
	)
}

func (r *AppRepository) InsertRepo(ctx context.Context, repo svcmodel.SCMRepo) (svcmodel.SCMRepo, error) {
	
	err := r.client.
		DB.From(r.entity).
		Update(model.ToDBSCMRepo(repo.Platform(), repo.RepoAbsURL(), repo.RepoOwnerIdentifier(), repo.RepoIdentifier())).
		Eq("id", repo.AppID().String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return r.ReadRepo(ctx, repo.AppID())
}

func (r *AppRepository) ReadRepo(ctx context.Context, appID uuid.UUID) (svcmodel.SCMRepo, error) {
	var resp model.SCMRepo
	err := r.client.
		DB.From(r.entity).
		Select("scm_repo").
		Single().
		Eq("id", appID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcSCMRepo(appID, resp.Data.Platform, resp.Data.RepoURL)
}

func (r *AppRepository) DeleteRepo(ctx context.Context, appID uuid.UUID) error {

	emptyRepo := svcmodel.NewEmptySCMRepo()
	err := r.client.
		DB.From(r.entity).
		Update(model.ToDBSCMRepo(emptyRepo.RepoIdentifier(), emptyRepo.RepoAbsURL(), emptyRepo.RepoOwnerIdentifier(), emptyRepo.Platform())).
		Eq("id", appID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return utils.WrapSupabaseDBErr(err)
	}

	return nil
}
