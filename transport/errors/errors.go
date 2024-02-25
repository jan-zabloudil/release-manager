package errors

import (
	"errors"
)

var (
	ErrInvalidBearer              = errors.New("invalid bearer token provided in authorization header")
	ErrAccessDeniedToAnonUser     = errors.New("access denied for anonymous user")
	ErrAccessDeniedToNonAdminUser = errors.New("access denied for non admin user")
	ErrHttpMethodNotAllowed       = errors.New("http method not allowed")
	ErrHttpNotFound               = errors.New("no route matched for requested uri")
)
