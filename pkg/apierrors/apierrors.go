package apierrors

import (
	"errors"
	"fmt"
)

var (
	errCodeUnauthorizedInvalidToken      = "ERR_UNAUTHORIZED_ACCESS_INVALID_TOKEN"
	errCodeForbiddenInsufficientUserRole = "ERR_FORBIDDEN_ACCESS_INSUFFICIENT_USER_ROLE"
	errCodeUserNotFound                  = "ERR_USER_NOT_FOUND"
	errCodeProjectNotFound               = "ERR_PROJECT_NOT_FOUND"
	errCodeEnvironmentNotFound           = "ERR_ENVIRONMENT_NOT_FOUND"
	errCodeProjectUnprocessable          = "ERR_PROJECT_UNPROCESSABLE"
	errCodeEnvironmentUnprocessable      = "ERR_ENVIRONMENT_UNPROCESSABLE"
	errCodeEnvironmentDuplicateName      = "ERR_ENVIRONMENT_DUPLICATE_NAME"
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

func (e *APIError) WithMessage(msg string) *APIError {
	e.Message = msg
	return e
}

func NewUserNotFoundError() *APIError {
	return &APIError{
		Code:    errCodeUserNotFound,
		Message: "User not found",
	}
}

func NewProjectNotFoundError() *APIError {
	return &APIError{
		Code:    errCodeProjectNotFound,
		Message: "Project not found",
	}
}

func NewEnvironmentNotFoundError() *APIError {
	return &APIError{
		Code:    errCodeEnvironmentNotFound,
		Message: "Environment not found",
	}
}

func NewEnvironmentDuplicateNameError() *APIError {
	return &APIError{
		Code:    errCodeEnvironmentDuplicateName,
		Message: "environment name is already in use",
	}
}

func NewProjectUnprocessableError() *APIError {
	return &APIError{
		Code:    errCodeProjectUnprocessable,
		Message: "Project unprocessable",
	}
}

func NewEnvironmentUnprocessableError() *APIError {
	return &APIError{
		Code:    errCodeEnvironmentUnprocessable,
		Message: "Environment unprocessable",
	}
}

func NewUnauthorizedError() *APIError {
	return &APIError{
		Code:    errCodeUnauthorizedInvalidToken,
		Message: "Unauthorized access, invalid or expired token provided.",
	}
}

func NewForbiddenInsufficientUserRoleError() *APIError {
	return &APIError{
		Code:    errCodeForbiddenInsufficientUserRole,
		Message: "Forbidden access, insufficient user role.",
	}
}

func IsErrorWithCode(err error, code string) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == code
	}

	return false
}

func IsNotFoundError(err error) bool {
	return IsErrorWithCode(err, errCodeUserNotFound) ||
		IsErrorWithCode(err, errCodeProjectNotFound) ||
		IsErrorWithCode(err, errCodeEnvironmentNotFound)
}

func IsUnprocessableModelError(err error) bool {
	return IsErrorWithCode(err, errCodeProjectUnprocessable) || IsErrorWithCode(err, errCodeEnvironmentUnprocessable)
}

func IsUnauthorizedError(err error) bool {
	return IsErrorWithCode(err, errCodeUnauthorizedInvalidToken)
}

func IsForbiddenError(err error) bool {
	return IsErrorWithCode(err, errCodeForbiddenInsufficientUserRole)
}

func IsConflictError(err error) bool {
	return IsErrorWithCode(err, errCodeEnvironmentDuplicateName)
}
