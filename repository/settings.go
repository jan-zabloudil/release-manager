package repository

import (
	"context"

	"release-manager/pkg/dberrors"
	"release-manager/repository/model"
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
		entity: "settings",
	}
}

func (r *SettingsRepository) Update(ctx context.Context, s svcmodel.Settings) error {
	data := model.ToSettings(s)

	err := r.client.
		DB.From(r.entity).
		Upsert(&data).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return dberrors.NewUnknownError().Wrap(err)
	}

	return nil
}

func (r *SettingsRepository) Read(ctx context.Context) (svcmodel.Settings, error) {
	var resp model.Settings
	err := r.client.
		DB.From(r.entity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Settings{}, dberrors.NewUnknownError().Wrap(err)
	}

	return model.ToSvcSettings(resp), nil
}
