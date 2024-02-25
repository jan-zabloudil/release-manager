package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

	"github.com/nedpals/supabase-go"
)

type SettingsRepository struct {
	client *supabase.Client
	entity string
}

func NewSettingsRepository(c *supabase.Client) *SettingsRepository {
	return &SettingsRepository{
		client: c,
		entity: "organization_settings",
	}
}

func (r *SettingsRepository) Set(ctx context.Context, s svcmodel.Settings) (svcmodel.Settings, error) {

	var resp []model.Setting
	err := r.client.
		DB.From(r.entity).
		Upsert(model.ToDBSettings(s.OrganizationName, s.SlackToken, s.GithubToken, s.DefaultReleaseMsg)).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Settings{}, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcSettings(resp), nil
}

func (r *SettingsRepository) Read(ctx context.Context) (svcmodel.Settings, error) {

	var resp []model.Setting
	err := r.client.
		DB.From(r.entity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Settings{}, utils.WrapSupabaseDBErr(err) // TODO dořešit, proč Supabase nevrací pořádně žádnou chybu
	}

	return model.ToSvcSettings(resp), nil
}
