package apierrors

import (
	"errors"
	"fmt"
)

var (
	unauthorizedInvalidTokenErrCode      = "ERR_UNAUTHORIZED_ACCESS_INVALID_TOKEN"
	forbiddenInsufficientUserRoleErrCode = "ERR_FORBIDDEN_ACCESS_INSUFFICIENT_USER_ROLE"
	userNotFoundErrCode                  = "ERR_USER_NOT_FOUND"
)

type APIError struct {
	Code    string
	Message string
	Err     error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Code: %s, error: %s", e.Code, e.Err)
}

func (e *APIError) Wrap(err error) *APIError {
	return &APIError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

func NewUserNotFoundError() *APIError {
	return &APIError{
		Code:    userNotFoundErrCode,
		Message: "User not found",
	}
}

func NewUnauthorizedError() *APIError {
	return &APIError{
		Code:    unauthorizedInvalidTokenErrCode,
		Message: "Unauthorized access, invalid or expired token provided.",
	}
}

func NewForbiddenInsufficientUserRoleError() *APIError {
	return &APIError{
		Code:    forbiddenInsufficientUserRoleErrCode,
		Message: "Forbidden access, insufficient user role.",
	}
}

func IsNotFoundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case userNotFoundErrCode:
			return true
		default:
			return false
		}
	}

	return false
}

func IsUnauthorizedError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == unauthorizedInvalidTokenErrCode
	}

	return false
}

func IsForbiddenError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == forbiddenInsufficientUserRoleErrCode
	}

	return false
}
