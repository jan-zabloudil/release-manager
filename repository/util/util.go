package util

import (
	"errors"

	"release-manager/pkg/autherrors"
	"release-manager/pkg/dberrors"

	postgrestgo "github.com/nedpals/postgrest-go/pkg"
	"github.com/nedpals/supabase-go"
)

const (
	postgresSingleRecordFetchErrorCode = "PGRST116"
)

const (
	supabaseNotFoundErrorCode     = 404
	supabaseUnauthorizedErrorCode = 401
)

func ToDBError(err error) *dberrors.DBError {
	var postgreErr *postgrestgo.RequestError
	if errors.As(err, &postgreErr) {
		if postgreErr.Code == postgresSingleRecordFetchErrorCode {
			return dberrors.NewNotFoundError().Wrap(err)
		}
	}

	return dberrors.NewUnknownError().Wrap(err)
}

func ToAuthError(err error) *autherrors.AuthError {
	var supabaseErr *supabase.ErrorResponse
	if errors.As(err, &supabaseErr) {
		switch supabaseErr.Code {
		case supabaseUnauthorizedErrorCode, supabaseNotFoundErrorCode:
			return autherrors.NewInvalidTokenError().Wrap(err)
		}
	}

	return autherrors.NewUnknownError().Wrap(err)
}
