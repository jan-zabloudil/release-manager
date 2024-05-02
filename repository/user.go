package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

const (
	userDBEntity = "users"
)

type UserRepository struct {
	client *supabase.Client
}

func NewUserRepository(client *supabase.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (r *UserRepository) Read(ctx context.Context, userID uuid.UUID) (svcmodel.User, error) {
	var resp model.User
	err := r.client.
		DB.From(userDBEntity).
		Select("*").Single().
		Eq("id", userID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.User{}, util.ToDBError(err)
	}

	return model.ToSvcUser(resp), nil
}

func (r *UserRepository) ReadAll(ctx context.Context) ([]svcmodel.User, error) {
	var resp []model.User
	err := r.client.
		DB.From(userDBEntity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	return model.ToSvcUsers(resp), nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.client.Admin.DeleteUser(ctx, id.String()); err != nil {
		return util.ToDBError(err)
	}

	return nil
}
