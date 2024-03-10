package errors

import "errors"

var (
	ErrTagAlreadyExists       = errors.New("github api: git tag already exists")
	ErrInvalidTargetCommitish = errors.New("github api: invalid target commitish provided")
	UnknownErr                = errors.New("github api: unknown error")
	ErrResourceNotFound       = errors.New("github api: resource not found")
	ErrUnauthenticated        = errors.New("github api: authentication not successful")
	ErrForbidden              = errors.New("github api: forbidden access")
)
