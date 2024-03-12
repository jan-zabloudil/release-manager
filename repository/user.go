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
	client         *supabase.Client
	getByEmailFunc string
}

func NewUserRepository(c *supabase.Client) *UserRepository {
	return &UserRepository{
		client:         c,
		getByEmailFunc: "get_user_by_email",
	}
}

func (r *UserRepository) ReadForToken(ctx context.Context, token string) (svcmodel.User, error) {
	res, err := r.client.Auth.User(ctx, token)
	if err != nil {
		return svcmodel.User{}, utils.WrapSupabaseAuthErr(err)
	}

	return model.ToSvcUser(
		res.ID,
		res.Email,
		res.AppMetadata["is_admin"],
		res.UserMetadata["name"],
		res.UserMetadata["picture"],
		res.CreatedAt,
		res.UpdatedAt,
	)
}

func (r *UserRepository) Read(ctx context.Context, ID uuid.UUID) (svcmodel.User, error) {
	res, err := r.client.Admin.GetUser(ctx, ID.String())
	if err != nil {
		return svcmodel.User{}, utils.WrapSupabaseAdminErr(err)
	}

	return model.ToSvcUser(
		res.ID,
		res.Email,
		res.AppMetaData["is_admin"],
		res.UserMetaData["name"],
		res.UserMetaData["picture"],
		res.CreatedAt,
		res.UpdatedAt,
	)
}

func (r *UserRepository) ReadByEmail(ctx context.Context, email string) (svcmodel.User, error) {
	var resp []supabase.User
	input := map[string]interface{}{
		"p_email": email,
	}

	// Supabase Auth.DB does not enable to query user by email, therefore postgres function must be used
	if err := r.client.DB.Rpc(r.getByEmailFunc, input).ExecuteWithContext(ctx, &resp); err != nil {
		return svcmodel.User{}, utils.WrapSupabaseDBErr(err)
	}
	if err := utils.ValidateSingleRecordFetchAfterReadOperation(resp); err != nil {
		return svcmodel.User{}, err
	}

	u := resp[0]
	return model.ToSvcUser(
		u.ID,
		u.Email,
		u.AppMetadata["is_admin"],
		u.UserMetadata["name"],
		u.UserMetadata["picture"],
		u.CreatedAt,
		u.UpdatedAt,
	)
}
