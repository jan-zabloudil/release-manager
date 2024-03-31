package responseerrors

import (
	"errors"
	"net/http"

	"release-manager/pkg/apierrors"
)

var (
	// #nosec G101 This is a constant error code, no security risk.
	notBearerTokenFormatErrorCode = "ERR_TOKEN_NOT_BEARER_FORMAT"
	missingBearerTokenErrCode     = "ERR_MISSING_BEARER_TOKEN"
	defaultNotFoundErrCode        = "ERR_NOT_FOUND"
	defaultForbiddenErrCode       = "ERR_FORBIDDEN"
	defaultUnauthorizedErrCode    = "ERR_UNAUTHORIZED"
	methodNotAllowedErrCode       = "ERR_METHOD_NOT_ALLOWED"
	invalidResourceIDErrCode      = "ERR_INVALID_RESOURCE_ID"
	unknownErrCode                = "ERR_UNKNOWN"
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

func NewNotFoundError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusNotFound,
		Code:       defaultNotFoundErrCode,
	}
}

func NewNotBearerTokenFormatError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnauthorized,
		Code:       notBearerTokenFormatErrorCode,
	}
}

func NewMissingBearerTokenError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnauthorized,
		Code:       missingBearerTokenErrCode,
	}
}

func NewServerError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusInternalServerError,
		Code:       unknownErrCode,
	}
}

func NewMethodNotAllowedError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusMethodNotAllowed,
		Code:       methodNotAllowedErrCode,
	}
}

func NewForbiddenError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusForbidden,
		Code:       defaultForbiddenErrCode,
	}
}

func NewUnauthorizedError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusUnauthorized,
		Code:       defaultUnauthorizedErrCode,
	}
}

func NewInvalidResourceIDError() *ResponseError {
	return &ResponseError{
		StatusCode: http.StatusNotFound,
		Code:       invalidResourceIDErrCode,
	}
}
