package errors

import (
	"errors"
	"net/http"

	"release-manager/pkg/apierrors"
)

var (
	// #nosec G101 This is a constant error code, no security risk.
	errCodeNotBearerTokenFormat = "ERR_TOKEN_NOT_BEARER_FORMAT"
	// #nosec G101 This is a constant error code, no security risk.
	errCodeExpiredOrInvalidToken      = "ERR_EXPIRED_OR_INVALID_TOKEN"
	errCodeMissingBearerToken         = "ERR_MISSING_BEARER_TOKEN"
	errCodeDefaultNotFound            = "ERR_NOT_FOUND"
	errCodeDefaultForbidden           = "ERR_FORBIDDEN"
	errCodeDefaultUnauthorized        = "ERR_UNAUTHORIZED"
	errCodeMethodNotAllowed           = "ERR_METHOD_NOT_ALLOWED"
	errCodeInvalidResourceID          = "ERR_INVALID_RESOURCE_ID"
	errCodeDefaultBadRequest          = "ERR_BAD_REQUEST"
	errCodeDefaultUnprocessableEntity = "ERR_UNPROCESSABLE_ENTITY"
	errCodeDefaultConflict            = "ERR_CONFLICT"
	errCodeUnknown                    = "ERR_UNKNOWN"
)

type Error struct {
	StatusCode int
	Message    string
	Code       string
	Err        error
}

func (r *Error) Wrap(err error) *Error {
	var apiErr *apierrors.APIError
	if errors.As(err, &apiErr) {
		return &Error{
			StatusCode: r.StatusCode,
			Message:    apiErr.Message,
			Code:       apiErr.Code,
			Err:        err,
		}
	}

	return &Error{
		StatusCode: r.StatusCode,
		Message:    r.Message,
		Code:       r.Code,
		Err:        err,
	}
}

func (r *Error) WithMessage(msg string) *Error {
	r.Message = msg
	return r
}

func NewNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Code:       errCodeDefaultNotFound,
	}
}

func NewNotBearerTokenFormatError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeNotBearerTokenFormat,
	}
}

func NewMissingBearerTokenError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeMissingBearerToken,
	}
}

func NewExpiredOrInvalidTokenError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeExpiredOrInvalidToken,
	}
}

func NewServerError() *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errCodeUnknown,
	}
}

func NewMethodNotAllowedError() *Error {
	return &Error{
		StatusCode: http.StatusMethodNotAllowed,
		Code:       errCodeMethodNotAllowed,
	}
}

func NewForbiddenError() *Error {
	return &Error{
		StatusCode: http.StatusForbidden,
		Code:       errCodeDefaultForbidden,
	}
}

func NewUnauthorizedError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeDefaultUnauthorized,
	}
}

func NewInvalidResourceIDError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Code:       errCodeInvalidResourceID,
	}
}

func NewBadRequestError() *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Code:       errCodeDefaultBadRequest,
	}
}

func NewUnprocessableEntityError() *Error {
	return &Error{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       errCodeDefaultUnprocessableEntity,
	}
}

func NewConflictError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Code:       errCodeDefaultConflict,
	}
}
