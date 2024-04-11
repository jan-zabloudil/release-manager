package autherrors

import (
	"errors"
	"fmt"
)

var (
	// #nosec G101 This is a constant error code, no security risk.
	invalidTokenErrCode  = "ERR_AUTH_INVALID_TOKEN"
	unknownErrCode       = "ERR_AUTH_UNKNOWN"
	invalidUserIDErrCode = "ERR_AUTH_INVALID_USER_ID"
)

type AuthError struct {
	Code string
	Err  error
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("Code: %s, error: %s", e.Code, e.Err)
}

func (e *AuthError) Wrap(err error) *AuthError {
	return &AuthError{
		Code: e.Code,
		Err:  err,
	}
}

func NewInvalidTokenError() *AuthError {
	return &AuthError{
		Code: invalidTokenErrCode,
	}
}

func NewUnknownError(err error) *AuthError {
	return &AuthError{
		Code: unknownErrCode,
		Err:  err,
	}
}

func NewInvalidUserIDError() *AuthError {
	return &AuthError{
		Code: invalidUserIDErrCode,
	}
}

func IsInvalidTokenError(err error) bool {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Code == invalidTokenErrCode
	}

	return false
}
