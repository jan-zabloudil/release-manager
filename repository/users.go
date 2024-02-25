package repository

import (
	"context"
	"errors"
	"fmt"

	reperr "github.com/jan-zabloudil/release-manager/repository/errors"
	"github.com/jan-zabloudil/release-manager/repository/model"
	svcerr "github.com/jan-zabloudil/release-manager/service/errors"
	svcmodel "github.com/jan-zabloudil/release-manager/service/model"
	"github.com/nedpals/supabase-go"
)

type UserRepository struct {
	client *supabase.Client
}

func (r *UserRepository) ReadForToken(ctx context.Context, token string) (svcmodel.User, error) {

	res, err := r.client.Auth.User(ctx, token)
	if err != nil {
		var supabaseErr *supabase.ErrorResponse
		if errors.As(err, &supabaseErr) && supabaseErr.Code == 401 {
			return svcmodel.User{}, fmt.Errorf("%w: %s", svcerr.ErrUserAuthenticationFailed, err.Error())
		}

		return svcmodel.User{}, fmt.Errorf("%w: %s", reperr.SupabaseGeneralErr, err.Error())
	}

	u, err := model.AuthToSvcUser(*res)
	if err != nil {
		return svcmodel.User{}, err
	}

	return u, nil
}
