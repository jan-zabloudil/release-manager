package repository

import (
	"context"
	"fmt"

	"release-manager/pkg/id"
	"release-manager/repository/helper"
	"release-manager/repository/model"
	"release-manager/repository/query"
	svcerrors "release-manager/service/errors"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueGitTagPerProjectConstraintName = "unique_git_tag_per_project"
)

type ReleaseRepository struct {
	dbpool             *pgxpool.Pool
	githubURLGenerator githubURLGenerator
	fileURLGenerator   fileURLGenerator
}

func NewReleaseRepository(
	pool *pgxpool.Pool,
	githubURLGenerator githubURLGenerator,
	fileURLGenerator fileURLGenerator,
) *ReleaseRepository {
	return &ReleaseRepository{
		dbpool:             pool,
		githubURLGenerator: githubURLGenerator,
		fileURLGenerator:   fileURLGenerator,
	}
}

func (r *ReleaseRepository) CreateRelease(ctx context.Context, rls svcmodel.Release) error {
	if _, err := r.dbpool.Exec(ctx, query.CreateRelease, pgx.NamedArgs{
		"id":           rls.ID,
		"projectID":    rls.ProjectID,
		"releaseTitle": rls.ReleaseTitle,
		"releaseNotes": rls.ReleaseNotes,
		"gitTagName":   rls.GitTagName,
		"createdBy":    rls.AuthorUserID,
		"createdAt":    rls.CreatedAt,
		"updatedAt":    rls.UpdatedAt,
	}); err != nil {
		if helper.IsUniqueConstraintViolation(err, uniqueGitTagPerProjectConstraintName) {
			return svcerrors.NewReleaseGitTagAlreadyUsedError().Wrap(err)
		}

		return err
	}

	return nil
}

func (r *ReleaseRepository) ReadRelease(ctx context.Context, releaseID id.Release) (svcmodel.Release, error) {
	return r.readRelease(ctx, r.dbpool, query.ReadRelease, pgx.NamedArgs{
		"releaseID": releaseID,
	})
}

func (r *ReleaseRepository) ReadReleaseForProject(ctx context.Context, projectID uuid.UUID, releaseID id.Release) (svcmodel.Release, error) {
	return r.readRelease(ctx, r.dbpool, query.ReadReleaseForProject, pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": releaseID,
	})
}

func (r *ReleaseRepository) UpdateRelease(
	ctx context.Context,
	releaseID id.Release,
	updateFn func(r svcmodel.Release) (svcmodel.Release, error),
) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		rls, err := r.readRelease(ctx, tx, query.AppendForUpdate(query.ReadRelease), pgx.NamedArgs{
			"releaseID": releaseID,
		})
		if err != nil {
			return fmt.Errorf("reading release: %w", err)
		}

		rls, err = updateFn(rls)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(ctx, query.UpdateRelease, pgx.NamedArgs{
			"releaseID":    rls.ID,
			"releaseTitle": rls.ReleaseTitle,
			"releaseNotes": rls.ReleaseNotes,
			"updatedAt":    rls.UpdatedAt,
		}); err != nil {
			return fmt.Errorf("updating release: %w", err)
		}

		return nil
	})
}

func (r *ReleaseRepository) DeleteReleaseByGitTag(ctx context.Context, githubOwnerSlug, githubRepoSlug, gitTag string) error {
	return r.deleteRelease(ctx, r.dbpool, query.DeleteReleaseByGitTag, pgx.NamedArgs{
		"ownerSlug":  githubOwnerSlug,
		"repoSlug":   githubRepoSlug,
		"gitTagName": gitTag,
	})
}

func (r *ReleaseRepository) DeleteRelease(ctx context.Context, releaseID id.Release) error {
	return r.deleteRelease(ctx, r.dbpool, query.DeleteRelease, pgx.NamedArgs{
		"releaseID": releaseID,
	})
}

func (r *ReleaseRepository) ListReleasesForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Release, error) {
	releases, err := helper.ListValues[model.Release](ctx, r.dbpool, query.ListReleasesForProject, pgx.NamedArgs{
		"projectID": projectID,
	})
	if err != nil {
		return nil, err
	}

	return model.ToSvcReleases(releases, r.githubURLGenerator.GenerateGitTagURL, r.fileURLGenerator.GenerateFileURL)
}

func (r *ReleaseRepository) CreateDeployment(ctx context.Context, dpl svcmodel.Deployment) error {
	if _, err := r.dbpool.Exec(ctx, query.CreateDeployment, pgx.NamedArgs{
		"id":            dpl.ID,
		"releaseID":     dpl.Release.ID,
		"environmentID": dpl.Environment.ID,
		"deployedBy":    dpl.DeployedByUserID,
		"deployedAt":    dpl.DeployedAt,
	}); err != nil {
		return err
	}

	return nil
}

func (r *ReleaseRepository) ListDeploymentsForProject(ctx context.Context, params svcmodel.ListDeploymentsFilterParams, projectID uuid.UUID) ([]svcmodel.Deployment, error) {
	listQuery := query.ListDeploymentsForProject
	if params.LatestOnly != nil && *params.LatestOnly {
		listQuery = query.AppendLimit(listQuery, 1)
	}

	// Release and Environment IDs are filter params that are optional and can be nil
	dpls, err := helper.ListValues[model.Deployment](ctx, r.dbpool, listQuery, pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": params.ReleaseID,
		"envID":     params.EnvironmentID,
	})
	if err != nil {
		return nil, err
	}

	return model.ToSvcDeployments(dpls)
}

func (r *ReleaseRepository) ReadLastDeploymentForRelease(ctx context.Context, releaseID id.Release) (svcmodel.Deployment, error) {
	dpl, err := helper.ReadValue[model.Deployment](ctx, r.dbpool, query.ReadLastDeploymentForRelease, pgx.NamedArgs{
		"releaseID": releaseID,
	})
	if err != nil {
		if helper.IsNotFound(err) {
			return svcmodel.Deployment{}, svcerrors.NewDeploymentNotFoundError().Wrap(err)
		}

		return svcmodel.Deployment{}, err
	}

	return model.ToSvcDeployment(dpl)
}

func (r *ReleaseRepository) readRelease(ctx context.Context, q helper.Querier, query string, args pgx.NamedArgs) (svcmodel.Release, error) {
	rls, err := helper.ReadValue[model.Release](ctx, q, query, args)
	if err != nil {
		if helper.IsNotFound(err) {
			return svcmodel.Release{}, svcerrors.NewReleaseNotFoundError().Wrap(err)
		}

		return svcmodel.Release{}, err
	}

	return model.ToSvcRelease(rls, r.githubURLGenerator.GenerateGitTagURL, r.fileURLGenerator.GenerateFileURL)
}

func (r *ReleaseRepository) deleteRelease(ctx context.Context, e helper.ExecExecutor, query string, args pgx.NamedArgs) error {
	result, err := e.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return svcerrors.NewReleaseNotFoundError()
	}

	return nil
}
