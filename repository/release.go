package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/nedpals/supabase-go"
)

const releaseDBEntity = "releases"

type ReleaseRepository struct {
	client *supabase.Client
}

func NewReleaseRepository(c *supabase.Client) *ReleaseRepository {
	return &ReleaseRepository{
		client: c,
	}
}

func (r *ReleaseRepository) Create(ctx context.Context, rls svcmodel.Release) error {
	data := model.ToRelease(rls)

	err := r.client.
		DB.From(releaseDBEntity).
		Insert(&data).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}
