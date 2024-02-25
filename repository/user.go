package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

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
