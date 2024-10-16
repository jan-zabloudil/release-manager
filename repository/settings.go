package repository

import (
	"context"
	"fmt"

	"release-manager/repository/helper"
	"release-manager/repository/model"
	"release-manager/repository/query"
	svcmodel "release-manager/service/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingsRepository struct {
	dbpool *pgxpool.Pool
}

func NewSettingsRepository(pool *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{
		dbpool: pool,
	}
}

func (r *SettingsRepository) Read(ctx context.Context) (svcmodel.Settings, error) {
	return r.read(ctx, r.dbpool, query.ReadSettings)
}

func (r *SettingsRepository) Upsert(
	ctx context.Context,
	fn svcmodel.UpdateSettingsFunc,
) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		s, err := r.read(ctx, tx, query.AppendForUpdate(query.ReadSettings))
		if err != nil {
			return fmt.Errorf("reading settings: %w", err)
		}

		s, err = fn(s)
		if err != nil {
			return err
		}

		sv, err := model.ToSettingsValues(s)
		if err != nil {
			return fmt.Errorf("converting service settigs model to repository model: %w", err)
		}

		// In current implementation, there are up to 4 settings values to upsert.
		// It is ok to update them by one in a loop.
		for _, v := range sv {
			if _, err := tx.Exec(ctx, query.UpsertSettings, pgx.NamedArgs{
				"key":   v.Key,
				"value": v.Value,
			}); err != nil {
				return fmt.Errorf("upserting setting value: %w", err)
			}
		}

		return nil
	})
}

func (r *SettingsRepository) read(ctx context.Context, q helper.Querier, query string) (svcmodel.Settings, error) {
	s, err := helper.ListValues[model.SettingsValue](ctx, q, query, nil)
	if err != nil {
		return svcmodel.Settings{}, err
	}

	return model.ToSvcSettings(s)
}
