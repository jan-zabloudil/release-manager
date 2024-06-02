package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/query"
	svcmodel "release-manager/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
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

func (r *DeploymentRepository) ListForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Deployment, error) {
	var dpls []model.Deployment

	if err := pgxscan.Select(ctx, r.dbpool, &dpls, query.ListDeploymentsForProject, pgx.NamedArgs{
		"projectID": projectID,
	}); err != nil {
		return nil, err
	}

	return model.ToSvcDeployments(dpls)
}
