package repository

import (
	"context"

	"release-manager/repository/query"
	svcmodel "release-manager/service/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeploymentRepository struct {
	dbpool *pgxpool.Pool
}

func NewDeploymentRepository(pool *pgxpool.Pool) *DeploymentRepository {
	return &DeploymentRepository{
		dbpool: pool,
	}
}

func (r *DeploymentRepository) Create(ctx context.Context, dpl svcmodel.Deployment) error {
	_, err := r.dbpool.Exec(ctx, query.CreateDeployment, pgx.NamedArgs{
		"id":            dpl.ID,
		"releaseID":     dpl.Release.ID,
		"environmentID": dpl.Environment.ID,
		"deployedBy":    dpl.DeployedByUserID,
		"deployedAt":    dpl.DeployedAt,
	})
	if err != nil {
		return err
	}

	return nil
}
