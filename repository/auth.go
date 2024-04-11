package repository

import (
	"context"

	"release-manager/pkg/autherrors"
	"release-manager/repository/util"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type AuthRepository struct {
	client *supabase.Client
}

func NewAuthRepository(client *supabase.Client) *AuthRepository {
	return &AuthRepository{client}
}

func (r *AuthRepository) ReadUserIDForToken(ctx context.Context, token string) (uuid.UUID, error) {
	user, err := r.client.Auth.User(ctx, token)
	if err != nil {
		return uuid.Nil, util.ToAuthError(err)
	}

	id, err := uuid.Parse(user.ID)
	if err != nil {
		return uuid.Nil, autherrors.NewInvalidUserIDError().Wrap(err)
	}

	return id, nil
}
