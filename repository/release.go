package repository

import (
	"context"
	"errors"

	"release-manager/pkg/apierrors"
	"release-manager/repository/model"
	"release-manager/repository/query"
	svcmodel "release-manager/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		return err
	}

	return nil
}

func (r *ReleaseRepository) ReadForProject(ctx context.Context, projectID, releaseID uuid.UUID) (svcmodel.Release, error) {
	var rls model.Release

	err := pgxscan.Get(ctx, r.dbpool, &rls, query.ReadReleaseForProject, pgx.NamedArgs{"projectID": projectID, "releaseID": releaseID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.Release{}, apierrors.NewReleaseNotFoundError().Wrap(err)
		}

		return svcmodel.Release{}, err
	}

	return model.ToSvcRelease(rls), nil
}
