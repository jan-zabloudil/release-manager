package repository

import (
	"context"
	"errors"
	"fmt"

	"release-manager/pkg/apierrors"
	"release-manager/repository/model"
	"release-manager/repository/query"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueReleaseTitlePerProjectConstraintName = "unique_release_title_per_project"
)

type ReleaseRepository struct {
	dbpool *pgxpool.Pool
}

func NewReleaseRepository(pool *pgxpool.Pool) *ReleaseRepository {
	return &ReleaseRepository{
		dbpool: pool,
	}
}

func (r *ReleaseRepository) Create(ctx context.Context, rls svcmodel.Release) error {
	_, err := r.dbpool.Exec(ctx, query.CreateRelease, pgx.NamedArgs{
		"id":           rls.ID,
		"projectID":    rls.ProjectID,
		"releaseTitle": rls.ReleaseTitle,
		"releaseNotes": rls.ReleaseNotes,
		"createdBy":    rls.AuthorUserID,
		"createdAt":    rls.CreatedAt,
		"updatedAt":    rls.UpdatedAt,
	})
	if err != nil {
		if util.IsUniqueConstraintViolation(err, uniqueReleaseTitlePerProjectConstraintName) {
			return apierrors.NewReleaseDuplicateTitleError().Wrap(err)
		}

		return err
	}

	return nil
}

func (r *ReleaseRepository) Read(ctx context.Context, projectID, releaseID uuid.UUID) (svcmodel.Release, error) {
	// Project ID is not needed in the query because releaseID is primary key
	// But it is added for security reasons
	// To make sure that the release belongs to the project that is passed from the service
	return r.read(ctx, r.dbpool, query.ReadRelease, pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": releaseID,
	})
}

func (r *ReleaseRepository) read(ctx context.Context, q querier, readQuery string, args pgx.NamedArgs) (svcmodel.Release, error) {
	var rls model.Release

	err := pgxscan.Get(ctx, q, &rls, readQuery, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.Release{}, apierrors.NewReleaseNotFoundError().Wrap(err)
		}

		return svcmodel.Release{}, err
	}

	return model.ToSvcRelease(rls), nil
}

func (r *ReleaseRepository) Update(
	ctx context.Context,
	projectID,
	releaseID uuid.UUID,
	fn svcmodel.UpdateReleaseFunc,
) (rls svcmodel.Release, err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	// Project ID is not needed in the query because releaseID is primary key
	// But it is added for security reasons
	// To make sure that the release belongs to the project that is passed from the service
	rls, err = r.read(ctx, tx, query.AppendForUpdate(query.ReadRelease), pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": releaseID,
	})
	if err != nil {
		return svcmodel.Release{}, fmt.Errorf("failed to read release: %w", err)
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
		return svcmodel.Release{}, fmt.Errorf("failed to update release: %w", err)
	}

	return rls, nil
}

func (r *ReleaseRepository) Delete(ctx context.Context, projectID, releaseID uuid.UUID) error {
	// Project ID is not needed in the query because releaseID is primary key
	// But it is added for security reasons
	// To make sure that the release belongs to the project that is passed from the service
	result, err := r.dbpool.Exec(ctx, query.DeleteRelease, pgx.NamedArgs{
		"projectID": projectID,
		"releaseID": releaseID,
	})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apierrors.NewReleaseNotFoundError()
	}

	return nil
}

func (r *ReleaseRepository) ListForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Release, error) {
	var rls []model.Release

	err := pgxscan.Select(ctx, r.dbpool, &rls, query.ListReleasesForProject, pgx.NamedArgs{
		"projectID": projectID,
	})
	if err != nil {
		return nil, err
	}

	return model.ToSvcReleases(rls), nil
}
