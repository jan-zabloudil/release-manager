package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type TemplateRepository struct {
	client *supabase.Client
	entity string
}

func NewTemplateRepository(c *supabase.Client) *TemplateRepository {
	return &TemplateRepository{
		client: c,
		entity: "templates",
	}
}

func (r *TemplateRepository) Insert(ctx context.Context, t svcmodel.Template) (svcmodel.Template, error) {

	var resp []model.Template
	err := r.client.
		DB.From(r.entity).
		Insert(model.ToDBTemplate(
			t.ID,
			t.Type.TemplateType(),
			model.ToDBReleaseMsg(t.ReleaseMsg.Title, t.ReleaseMsg.Text, t.ReleaseMsg.Includes)),
		).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Template{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetch(resp); err != nil {
		return svcmodel.Template{}, err
	}

	return model.ToSvcTemplate(
		resp[0].ID,
		resp[0].Type,
		model.ToSvcReleaseMsg(resp[0].ReleaseMsg.Title, resp[0].ReleaseMsg.Text, resp[0].ReleaseMsg.Includes),
	)
}

func (r *TemplateRepository) ReadAll(ctx context.Context) ([]svcmodel.Template, error) {

	var resp []model.Template
	err := r.client.
		DB.From(r.entity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcTemplates(resp)
}

func (r *TemplateRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.Template, error) {

	var resp model.Template
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Template{}, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcTemplate(
		resp.ID,
		resp.Type,
		model.ToSvcReleaseMsg(resp.ReleaseMsg.Title, resp.ReleaseMsg.Text, resp.ReleaseMsg.Includes),
	)
}

func (r *TemplateRepository) Update(ctx context.Context, t svcmodel.Template) (svcmodel.Template, error) {

	var resp []model.Template
	err := r.client.
		DB.From(r.entity).
		Update(model.ToDBTemplate(
			t.ID,
			t.Type.TemplateType(),
			model.ToDBReleaseMsg(t.ReleaseMsg.Title, t.ReleaseMsg.Text, t.ReleaseMsg.Includes)),
		).
		Eq("id", t.ID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Template{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetch(resp); err != nil {
		return svcmodel.Template{}, err
	}

	return model.ToSvcTemplate(
		resp[0].ID,
		resp[0].Type,
		model.ToSvcReleaseMsg(resp[0].ReleaseMsg.Title, resp[0].ReleaseMsg.Text, resp[0].ReleaseMsg.Includes),
	)
}

func (r *TemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {

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
