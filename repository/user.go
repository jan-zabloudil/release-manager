package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type UserRepository struct {
	client *supabase.Client
}

func (r *UserRepository) ReadForToken(ctx context.Context, token string) (svcmodel.User, error) {
	res, err := r.client.Auth.User(ctx, token)
	if err != nil {
		return svcmodel.User{}, utils.WrapSupabaseAuthErr(err)
	}

	return model.ToSvcUser(
		res.ID,
		res.Email,
		res.AppMetadata["role"],
		res.UserMetadata["name"],
		res.UserMetadata["picture"],
		res.CreatedAt,
		res.UpdatedAt,
	)
}

func (r *UserRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.User, error) {
	res, err := r.client.Admin.GetUser(ctx, id.String())
	if err != nil {
		return svcmodel.User{}, utils.WrapSupabaseAdminErr(err)
	}

	return model.ToSvcUser(
		res.ID,
		res.Email,
		res.AppMetaData["role"],
		res.UserMetaData["name"],
		res.UserMetaData["picture"],
		res.CreatedAt,
		res.UpdatedAt,
	)
}

func (r *UserRepository) ReadAll(ctx context.Context) ([]svcmodel.User, error) {
	res, err := r.client.Admin.GetUsers(ctx)
	if err != nil {
		return nil, utils.WrapSupabaseAdminErr(err)
	}

	return model.ToSvcUsers(res)
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.client.Admin.DeleteUser(ctx, id.String()); err != nil {
		return utils.WrapSupabaseAdminErr(err)
	}

	return nil
}
