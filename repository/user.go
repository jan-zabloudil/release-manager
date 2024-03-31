package repository

import (
	"context"

	"release-manager/pkg/dberrors"
	"release-manager/repository/model"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type UserRepository struct {
	client *supabase.Client
	entity string
}

func NewUserRepository(client *supabase.Client) *UserRepository {
	return &UserRepository{
		client: client,
		entity: "users",
	}
}

func (r *UserRepository) Read(ctx context.Context, userID uuid.UUID) (svcmodel.User, error) {
	var resp model.User
	err := r.client.
		DB.From(r.entity).
		Select("*").Single().
		Eq("id", userID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.User{}, util.ToDBError(err)
	}

	u, err := svcmodel.ToUser(
		resp.ID,
		resp.Email,
		resp.Name,
		resp.AvatarURL,
		resp.Role,
		resp.CreatedAt.Time,
		resp.UpdatedAt.Time,
	)
	if err != nil {
		return svcmodel.User{}, dberrors.NewToSvcModelError().Wrap(err)
	}

	return u, nil
}

func (r *UserRepository) ReadAll(ctx context.Context) ([]svcmodel.User, error) {
	var resp []model.User
	err := r.client.
		DB.From(r.entity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	u, err := model.ToSvcUsers(resp)
	if err != nil {
		return nil, dberrors.NewToSvcModelError().Wrap(err)
	}

	return u, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.client.Admin.DeleteUser(ctx, id.String()); err != nil {
		return util.ToDBError(err)
	}

	return nil
}
