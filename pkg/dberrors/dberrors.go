package dberrors

import (
	"errors"
	"fmt"
)

var (
	errCodeNotFound            = "ERR_DB_RESOURCE_NOT_FOUND"
	errCodeCannotMapToSvcModel = "ERR_DB_CANNOT_MAP_TO_SVC_MODEL"
	errCodeUnknown             = "ERR_DB_UNKNOWN"
)

type DBError struct {
	Code string
	Err  error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("Code: %s, error: %s", e.Code, e.Err)
}

func (e *DBError) Wrap(err error) *DBError {
	return &DBError{
		Code: e.Code,
		Err:  err,
	}
}

func NewNotFoundError() *DBError {
	return &DBError{
		Code: errCodeNotFound,
	}
}

func NewToSvcModelError() *DBError {
	return &DBError{
		Code: errCodeCannotMapToSvcModel,
	}
}

func NewUnknownError() *DBError {
	return &DBError{
		Code: errCodeUnknown,
	}
}

func IsNotFoundError(err error) bool {
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		return dbErr.Code == errCodeNotFound
	}

	return false
}
