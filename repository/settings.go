package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/query"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
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
	return r.read(ctx, r.dbpool)
}

func (r *SettingsRepository) read(ctx context.Context, q pgxscan.Querier) (svcmodel.Settings, error) {
	var sv []model.SettingsValue

	err := pgxscan.Select(ctx, q, &sv, query.ReadSettings)
	if err != nil {
		return svcmodel.Settings{}, err
	}

	return model.ToSvcSettings(sv)
}

func (r *SettingsRepository) Update(
	ctx context.Context,
	fn svcmodel.UpdateSettingsFunc,
) (s svcmodel.Settings, err error) {
	err = util.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		s, err = r.read(ctx, tx)
		if err != nil {
			return err
		}

		s, err = fn(s)
		if err != nil {
			return err
		}

		if err := r.update(ctx, tx, s); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return svcmodel.Settings{}, err
	}

	return s, nil
}

// update settings are saved as key-value pairs in the database
// if key-value pair already exists, it is updated
// if key-value pair does not exist, it is inserted
func (r *SettingsRepository) update(ctx context.Context, tx pgx.Tx, s svcmodel.Settings) error {
	sv, err := model.ToSettingsValues(s)
	if err != nil {
		return err
	}

	// sending batch is implicitly transactional
	// https://github.com/jackc/pgx/issues/879
	batch := &pgx.Batch{}
	for _, v := range sv {
		batch.Queue(query.UpsertSettings, pgx.NamedArgs{"key": v.Key, "value": v.Value})
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range sv {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}
