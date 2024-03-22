package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type ReleaseRepository struct {
	client *supabase.Client
	entity string
}

func NewReleaseRepository(c *supabase.Client) *ReleaseRepository {
	return &ReleaseRepository{
		client: c,
		entity: "releases",
	}
}

func (r *ReleaseRepository) Insert(ctx context.Context, release svcmodel.Release) (svcmodel.Release, error) {

	var resp []model.Release
	err := r.client.DB.
		From(r.entity).
		Insert(model.ToDBRelease(
			release.ID,
			release.AppID,
			release.CreatedByUserID,
			release.SourceCode.Tag(),
			release.SourceCode.TargetCommitIsh(),
			release.Deployments.Dev,
			release.Deployments.Stg,
			release.Deployments.Prd,
			release.Title,
			release.ChangeLog,
			release.CreatedAt,
			release.UpdatedAt,
		)).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Release{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetchAfterWriteOperation(resp); err != nil {
		return svcmodel.Release{}, err
	}

	return model.ToSvcRelease(
		resp[0].ID,
		resp[0].AppID,
		resp[0].CreatedByUserID,
		resp[0].SourceCode.Tag,
		resp[0].SourceCode.TargetCommitIsh,
		resp[0].Deployments.Dev,
		resp[0].Deployments.Stg,
		resp[0].Deployments.Prd,
		resp[0].Title,
		resp[0].ChangeLog,
		resp[0].CreatedAt,
		resp[0].UpdatedAt,
	)
}

func (r *ReleaseRepository) Update(ctx context.Context, release svcmodel.Release) (svcmodel.Release, error) {

	var resp []model.Release
	err := r.client.DB.
		From(r.entity).
		Update(model.ToDBRelease(
			release.ID,
			release.AppID,
			release.CreatedByUserID,
			release.SourceCode.Tag(),
			release.SourceCode.TargetCommitIsh(),
			release.Deployments.Dev,
			release.Deployments.Stg,
			release.Deployments.Prd,
			release.Title,
			release.ChangeLog,
			release.CreatedAt,
			release.UpdatedAt,
		)).
		Eq("id", release.ID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Release{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetchAfterWriteOperation(resp); err != nil {
		return svcmodel.Release{}, err
	}

	return model.ToSvcRelease(
		resp[0].ID,
		resp[0].AppID,
		resp[0].CreatedByUserID,
		resp[0].SourceCode.Tag,
		resp[0].SourceCode.TargetCommitIsh,
		resp[0].Deployments.Dev,
		resp[0].Deployments.Stg,
		resp[0].Deployments.Prd,
		resp[0].Title,
		resp[0].ChangeLog,
		resp[0].CreatedAt,
		resp[0].UpdatedAt,
	)
}

func (r *ReleaseRepository) ReadAllForApp(ctx context.Context, appID uuid.UUID) ([]svcmodel.Release, error) {

	var resp []model.Release
	err := r.client.
		DB.From(r.entity).
		Select("*").
		Eq("app_id", appID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcReleases(resp)
}

func (r *ReleaseRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.Release, error) {

	var resp model.Release
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Release{}, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcRelease(
		resp.ID,
		resp.AppID,
		resp.CreatedByUserID,
		resp.SourceCode.Tag,
		resp.SourceCode.TargetCommitIsh,
		resp.Deployments.Dev,
		resp.Deployments.Stg,
		resp.Deployments.Prd,
		resp.Title,
		resp.ChangeLog,
		resp.CreatedAt,
		resp.UpdatedAt,
	)
}

func (r *ReleaseRepository) Delete(ctx context.Context, id uuid.UUID) error {

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
