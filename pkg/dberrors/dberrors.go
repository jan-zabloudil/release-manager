package dberrors

import (
	"errors"
	"fmt"
)

var (
	notFoundErrCode            = "ERR_DB_RESOURCE_NOT_FOUND"
	cannotMapToSvcModelErrCode = "ERR_DB_CANNOT_MAP_TO_SVC_MODEL"
	multipleOrNoRecordsErrCode = "ERR_DB_MULTIPLE_OR_NO_RECORDS_RETURNED"
	unknownErrCode             = "ERR_DB_UNKNOWN"
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
		Code: notFoundErrCode,
	}
}

func NewToSvcModelError() *DBError {
	return &DBError{
		Code: cannotMapToSvcModelErrCode,
	}
}

func NewMultipleOrNoRecordsReturnedError() *DBError {
	return &DBError{
		Code: multipleOrNoRecordsErrCode,
	}
}

func NewUnknownError(err error) *DBError {
	return &DBError{
		Code: unknownErrCode,
		Err:  err,
	}
}

func IsNotFoundError(err error) bool {
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		return dbErr.Code == notFoundErrCode
	}

	return false
}
