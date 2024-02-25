package errors

import "errors"

var ErrUserAuthenticationFailed = errors.New("user authentication by bearer token unsuccessful")
