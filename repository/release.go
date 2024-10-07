package repository

import (
	"context"
	"errors"
	"fmt"

	"release-manager/repository/model"
	"release-manager/repository/query"
	"release-manager/repository/util"
	svcerrors "release-manager/service/errors"
	svcmodel "release-manager/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
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
	_, err := r.dbpool.Exec(ctx, query.CreateRelease, pgx.NamedArgs{
		"id":           rls.ID,
		"projectID":    rls.ProjectID,
		"releaseTitle": rls.ReleaseTitle,
		"releaseNotes": rls.ReleaseNotes,
		"gitTagName":   rls.GitTagName,
		"createdBy":    rls.AuthorUserID,
		"createdAt":    rls.CreatedAt,
		"updatedAt":    rls.UpdatedAt,
	})
	if err != nil {
		if util.IsUniqueConstraintViolation(err, uniqueGitTagPerProjectConstraintName) {
			return svcerrors.NewReleaseGitTagAlreadyUsedError().Wrap(err)
		}

		return err
	}

	return nil
}

func (r *ReleaseRepository) ReadRelease(ctx context.Context, releaseID uuid.UUID) (svcmodel.Release, error) {
	return r.readRelease(ctx, r.dbpool, query.ReadRelease, pgx.NamedArgs{"releaseID": releaseID})
}

func (r *ReleaseRepository) ReadReleaseForProject(ctx context.Context, projectID, releaseID uuid.UUID) (svcmodel.Release, error) {
	// Project ID is not needed in the query because releaseID is primary key
	// But it is added for security reasons
	// To make sure that the release belongs to the project that is passed from the service
	return r.readRelease(ctx, r.dbpool, query.ReadRelease, pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": releaseID,
	})
}

func (r *ReleaseRepository) readRelease(ctx context.Context, q querier, readQuery string, args pgx.NamedArgs) (svcmodel.Release, error) {
	var rls model.Release

	err := pgxscan.Get(ctx, q, &rls, readQuery, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.Release{}, svcerrors.NewReleaseNotFoundError().Wrap(err)
		}

		return svcmodel.Release{}, err
	}

	return model.ToSvcRelease(rls, r.githubURLGenerator.GenerateGitTagURL, r.fileURLGenerator.GenerateFileURL)
}

func (r *ReleaseRepository) UpdateRelease(
	ctx context.Context,
	projectID,
	releaseID uuid.UUID,
	fn svcmodel.UpdateReleaseFunc,
) (rls svcmodel.Release, err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	// Project ID is not needed in the query because releaseID is primary key
	// But it is added for security reasons
	// To make sure that the release belongs to the project that is passed from the service
	rls, err = r.readRelease(ctx, tx, query.AppendForUpdate(query.ReadRelease), pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": releaseID,
	})
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("reading release: %w", err)
	}

	// Update the release
	rls, err = fn(rls)
	if err != nil {
		return svcmodel.Release{}, err
	}

	_, err = tx.Exec(ctx, query.UpdateRelease, pgx.NamedArgs{
		"releaseID":    rls.ID,
		"releaseTitle": rls.ReleaseTitle,
		"releaseNotes": rls.ReleaseNotes,
		"updatedAt":    rls.UpdatedAt,
	})
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("updating release: %w", err)
	}

	return rls, nil
}

func (r *ReleaseRepository) DeleteReleaseByGitTag(ctx context.Context, githubOwnerSlug, githubRepoSlug, gitTag string) error {
	return r.deleteRelease(ctx, query.DeleteReleaseByGitTag, pgx.NamedArgs{
		"ownerSlug":  githubOwnerSlug,
		"repoSlug":   githubRepoSlug,
		"gitTagName": gitTag,
	})
}

func (r *ReleaseRepository) DeleteRelease(ctx context.Context, projectID, releaseID uuid.UUID) error {
	// Project ID is not needed in the query because releaseID is primary key
	// But it is added for security reasons
	// To make sure that the release belongs to the project that is passed from the service
	return r.deleteRelease(ctx, query.DeleteRelease, pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": releaseID,
	})
}

func (r *ReleaseRepository) deleteRelease(ctx context.Context, deleteQuery string, args pgx.NamedArgs) error {
	result, err := r.dbpool.Exec(ctx, deleteQuery, args)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return svcerrors.NewReleaseNotFoundError()
	}

	return nil
}

func (r *ReleaseRepository) ListReleasesForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Release, error) {
	var dbReleases []model.Release

	err := pgxscan.Select(ctx, r.dbpool, &dbReleases, query.ListReleasesForProject, pgx.NamedArgs{
		"projectID": projectID,
	})
	if err != nil {
		return nil, err
	}

	return model.ToSvcReleases(dbReleases, r.githubURLGenerator.GenerateGitTagURL, r.fileURLGenerator.GenerateFileURL)
}

func (r *ReleaseRepository) CreateDeployment(ctx context.Context, dpl svcmodel.Deployment) error {
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

func (r *ReleaseRepository) ListDeploymentsForProject(ctx context.Context, params svcmodel.DeploymentFilterParams, projectID uuid.UUID) ([]svcmodel.Deployment, error) {
	var dpls []model.Deployment

	listQuery := query.ListDeploymentsForProject
	if params.LatestOnly != nil && *params.LatestOnly {
		listQuery = query.AppendLimit(listQuery, 1)
	}

	if err := pgxscan.Select(ctx, r.dbpool, &dpls, listQuery, pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": params.ReleaseID,
		"envID":     params.EnvironmentID,
	}); err != nil {
		return nil, err
	}

	return model.ToSvcDeployments(dpls)
}

func (r *ReleaseRepository) ReadLastDeploymentForRelease(ctx context.Context, projectID, releaseID uuid.UUID) (svcmodel.Deployment, error) {
	var dpl model.Deployment

	// Project ID is not needed in the query because releaseID is primary key
	// But it is added for security reasons
	// To make sure that the release (and therefore deployment) belongs to the project that is passed from the service
	err := pgxscan.Get(ctx, r.dbpool, &dpl, query.ReadLastDeploymentForRelease, pgx.NamedArgs{
		"releaseID": releaseID,
		"projectID": projectID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.Deployment{}, svcerrors.NewDeploymentNotFoundError().Wrap(err)
		}

		return svcmodel.Deployment{}, err
	}

	return model.ToSvcDeployment(dpl)
}
