package responseerrors

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

type ResponseError struct {
	StatusCode int
	Message    string
	Code       string
	Err        error
}

func (r *ResponseError) Wrap(err error) *ResponseError {
	var apiErr *apierrors.APIError
	if errors.As(err, &apiErr) {
		return &ResponseError{
			StatusCode: r.StatusCode,
			Message:    apiErr.Message,
			Code:       apiErr.Code,
			Err:        err,
		}
	}

	return &ResponseError{
		StatusCode: r.StatusCode,
		Message:    r.Message,
		Code:       r.Code,
		Err:        err,
	}
}

func (r *ResponseError) WithMessage(msg string) *ResponseError {
	r.Message = msg
	return r
}

func NewNotFoundError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusNotFound,
		Code:       errCodeDefaultNotFound,
	}
}

func NewNotBearerTokenFormatError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeNotBearerTokenFormat,
	}
}

func NewMissingBearerTokenError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeMissingBearerToken,
	}
}

func NewExpiredOrInvalidTokenError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeExpiredOrInvalidToken,
	}
}

func NewServerError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusInternalServerError,
		Code:       errCodeUnknown,
	}
}

func NewMethodNotAllowedError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusMethodNotAllowed,
		Code:       errCodeMethodNotAllowed,
	}
}

func NewForbiddenError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusForbidden,
		Code:       errCodeDefaultForbidden,
	}
}

func NewUnauthorizedError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeDefaultUnauthorized,
	}
}

func NewInvalidResourceIDError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusNotFound,
		Code:       errCodeInvalidResourceID,
	}
}

func NewBadRequestError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusBadRequest,
		Code:       errCodeDefaultBadRequest,
	}
}

func NewUnprocessableEntityError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       errCodeDefaultUnprocessableEntity,
	}
}

func NewConflictError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusConflict,
		Code:       errCodeDefaultConflict,
	}
}
