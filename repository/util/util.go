package util

import (
	"errors"

	"release-manager/pkg/dberrors"

	postgrestgo "github.com/nedpals/postgrest-go/pkg"
)

const (
	postgresSingleRecordFetchErrorCode = "PGRST116"
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
