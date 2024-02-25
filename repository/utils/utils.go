package utils

import (
	"errors"
	"fmt"

	reperr "release-manager/repository/errors"

	postgrestgo "github.com/nedpals/postgrest-go/pkg"
	"github.com/nedpals/supabase-go"
)

func ValidateSingleRecordFetch[T any](s []T) error {
	if len(s) != 1 {
		return reperr.ErrMultipleOrNoRecordsReturnedAfterWriteOperation
	}

	return nil
}

func WrapSupabaseDBErr(err error) error {
	var postgreErr *postgrestgo.RequestError

	if errors.As(err, &postgreErr) {
		switch postgreErr.Code {
		case "PGRST116":
			return fmt.Errorf("%w: %s", reperr.ErrResourceNotFound, err.Error())
		case "42P01":
			return fmt.Errorf("%w: %s", reperr.ErrEntityNotExists, err.Error())
		case "42501":
			return fmt.Errorf("%w: %s", reperr.ErrAccessToResourceDenied, err.Error())
		}
	}

	return fmt.Errorf("%w: %s", reperr.UnknownErr, err.Error())
}

func WrapSupabaseAuthErr(err error) error {
	var supabaseErr *supabase.ErrorResponse
	if errors.As(err, &supabaseErr) {
		switch supabaseErr.Code {
		case 401:
			return fmt.Errorf("%w: %s", reperr.ErrUserAuthenticationFailed, err.Error())
		case 404:
			return fmt.Errorf("%w: %s", reperr.ErrUserAuthenticationFailed, err.Error())
		}
	}

	return fmt.Errorf("%w: %s", reperr.UnknownErr, err.Error())
}

func WrapSupabaseAdminErr(err error) error {
	var supabaseErr *supabase.ErrorResponse
	if errors.As(err, &supabaseErr) {
		switch supabaseErr.Code {
		case 401:
			return fmt.Errorf("%w: %s", reperr.ErrAccessToResourceDenied, err.Error())
		case 404:
			return fmt.Errorf("%w: %s", reperr.ErrResourceNotFound, err.Error())
		}
	}

	return fmt.Errorf("%w: %s", reperr.UnknownErr, err.Error())
}
