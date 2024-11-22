package errors

import (
	"errors"
	"net/http"

	svcerrors "release-manager/service/errors"
)

var (
	// #nosec G101 This is a constant error code, no security risk.
	errCodeNotBearerTokenFormat = "ERR_TOKEN_NOT_BEARER_FORMAT"
	// #nosec G101 This is a constant error code, no security risk.
	errCodeExpiredOrInvalidToken   = "ERR_EXPIRED_OR_INVALID_TOKEN"
	errCodeMissingBearerToken      = "ERR_MISSING_BEARER_TOKEN"
	errCodeDefaultNotFound         = "ERR_NOT_FOUND"
	errCodeDefaultForbidden        = "ERR_FORBIDDEN"
	errCodeDefaultUnauthorized     = "ERR_UNAUTHORIZED"
	errCodeDefaultMethodNotAllowed = "ERR_METHOD_NOT_ALLOWED"
	errCodeDefaultBadRequest       = "ERR_BAD_REQUEST"
	errCodeDefaultConflict         = "ERR_CONFLICT"
	errCodeInvalidRequestPayload   = "ERR_INVALID_REQUEST_PAYLOAD"
	errCodeInvalidURLParams        = "ERR_INVALID_URL_PARAMS"
	errCodeUnknown                 = "ERR_UNKNOWN"
)

type Error struct {
	StatusCode int
	Message    string
	Code       string
	Err        error
	Data       any
}

func (r *Error) Wrap(err error) *Error {
	var svcErr *svcerrors.Error
	if errors.As(err, &svcErr) {
		return &Error{
			StatusCode: r.StatusCode,
			Message:    svcErr.Message,
			Code:       svcErr.Code,
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

func (r *Error) WithData(data any) *Error {
	r.Data = data
	return r
}

func NewDefaultNotFoundError() *Error {
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

func NewUnknownError() *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errCodeUnknown,
	}
}

func NewDefaultMethodNotAllowedError() *Error {
	return &Error{
		StatusCode: http.StatusMethodNotAllowed,
		Code:       errCodeDefaultMethodNotAllowed,
	}
}

func NewDefaultForbiddenError() *Error {
	return &Error{
		StatusCode: http.StatusForbidden,
		Code:       errCodeDefaultForbidden,
	}
}

func NewDefaultUnauthorizedError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Code:       errCodeDefaultUnauthorized,
	}
}

func NewDefaultBadRequestError() *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Code:       errCodeDefaultBadRequest,
	}
}

func NewInvalidRequestPayloadError() *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Code:       errCodeInvalidRequestPayload,
	}
}

func NewInvalidURLParamsError() *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Code:       errCodeInvalidURLParams,
	}
}

func NewDefaultConflictError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Code:       errCodeDefaultConflict,
	}
}
