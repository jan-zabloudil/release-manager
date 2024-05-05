package apierrors

import (
	"errors"
	"fmt"
)

var (
	errCodeUnauthorizedInvalidToken                = "ERR_UNAUTHORIZED_ACCESS_INVALID_TOKEN"
	errCodeForbiddenInsufficientUserRole           = "ERR_FORBIDDEN_ACCESS_INSUFFICIENT_USER_ROLE"
	errCodeForbiddenInsufficientProjectRole        = "ERR_FORBIDDEN_ACCESS_INSUFFICIENT_PROJECT_ROLE"
	errCodeUserNotFound                            = "ERR_USER_NOT_FOUND"
	errCodeProjectNotFound                         = "ERR_PROJECT_NOT_FOUND"
	errCodeEnvironmentNotFound                     = "ERR_ENVIRONMENT_NOT_FOUND"
	errCodeProjectUnprocessable                    = "ERR_PROJECT_UNPROCESSABLE"
	errCodeEnvironmentUnprocessable                = "ERR_ENVIRONMENT_UNPROCESSABLE"
	errCodeEnvironmentDuplicateName                = "ERR_ENVIRONMENT_DUPLICATE_NAME"
	errCodeSettingsUnprocessable                   = "ERR_SETTINGS_UNPROCESSABLE"
	errCodeProjectInvitationUnprocessable          = "ERR_PROJECT_INVITATION_UNPROCESSABLE"
	errCodeProjectInvitationAlreadyExists          = "ERR_PROJECT_INVITATION_ALREADY_EXISTS"
	errCodeProjectInvitationNotFound               = "ERR_PROJECT_INVITATION_NOT_FOUND"
	errCodeProjectMemberAlreadyExists              = "ERR_PROJECT_MEMBER_ALREADY_EXISTS"
	errCodeGithubIntegrationNotEnabled             = "ERR_GITHUB_INTEGRATION_NOT_ENABLED"
	errCodeGithubIntegrationUnauthorized           = "ERR_GITHUB_INTEGRATION_UNAUTHORIZED"
	errCodeGithubIntegrationForbidden              = "ERR_GITHUB_INTEGRATION_FORBIDDEN"
	errCodeGithubRepositoryNotConfiguredForProject = "ERR_GITHUB_REPOSITORY_NOT_CONFIGURED_FOR_PROJECT"
	errCodeGithubRepositoryNotFound                = "ERR_GITHUB_REPOSITORY_NOT_FOUND"
	errCodeProjectMemberNotFound                   = "ERR_PROJECT_MEMBER_NOT_FOUND"
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

func NewSettingsUnprocessableError() *APIError {
	return &APIError{
		Code:    errCodeSettingsUnprocessable,
		Message: "Settings unprocessable",
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

func NewForbiddenInsufficientProjectRoleError() *APIError {
	return &APIError{
		Code:    errCodeForbiddenInsufficientUserRole,
		Message: "Forbidden access, insufficient project role.",
	}
}

func NewProjectInvitationUnprocessableError() *APIError {
	return &APIError{
		Code:    errCodeProjectInvitationUnprocessable,
		Message: "Project invitation unprocessable",
	}
}

func NewProjectInvitationAlreadyExistsError() *APIError {
	return &APIError{
		Code:    errCodeProjectInvitationAlreadyExists,
		Message: "Project invitation already exists",
	}
}

func NewProjectInvitationNotFoundError() *APIError {
	return &APIError{
		Code:    errCodeProjectInvitationNotFound,
		Message: "Project invitation not found",
	}
}

func NewProjectMemberAlreadyExistsError() *APIError {
	return &APIError{
		Code:    errCodeProjectMemberAlreadyExists,
		Message: "Project member already exists",
	}
}

func NewGithubIntegrationNotEnabledError() *APIError {
	return &APIError{
		Code:    errCodeGithubIntegrationNotEnabled,
		Message: "Github integration is not enabled.",
	}
}

func NewGithubRepositoryNotConfiguredForProjectError() *APIError {
	return &APIError{
		Code:    errCodeGithubRepositoryNotConfiguredForProject,
		Message: "Github repository is not configured for the project.",
	}
}

func NewGithubIntegrationUnauthorizedError() *APIError {
	return &APIError{
		Code:    errCodeGithubIntegrationUnauthorized,
		Message: "Cannot access Github API, invalid or expired token.",
	}
}

func NewGithubIntegrationForbiddenError() *APIError {
	return &APIError{
		Code:    errCodeGithubIntegrationForbidden,
		Message: "Cannot access given resource via Github API.",
	}
}

func NewGithubRepositoryNotFoundError() *APIError {
	return &APIError{
		Code:    errCodeGithubRepositoryNotFound,
		Message: "Github repository not found among accessible repositories.",
	}
}

func NewProjectMemberNotFoundError() *APIError {
	return &APIError{
		Code:    errCodeProjectMemberNotFound,
		Message: "Project member not found",
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
		IsErrorWithCode(err, errCodeEnvironmentNotFound) ||
		IsErrorWithCode(err, errCodeProjectInvitationNotFound) ||
		IsErrorWithCode(err, errCodeGithubRepositoryNotFound) ||
		IsErrorWithCode(err, errCodeGithubRepositoryNotConfiguredForProject) ||
		IsErrorWithCode(err, errCodeGithubIntegrationNotEnabled) ||
		IsErrorWithCode(err, errCodeProjectMemberNotFound)
}

func IsUnprocessableModelError(err error) bool {
	return IsErrorWithCode(err, errCodeProjectUnprocessable) ||
		IsErrorWithCode(err, errCodeEnvironmentUnprocessable) ||
		IsErrorWithCode(err, errCodeSettingsUnprocessable) ||
		IsErrorWithCode(err, errCodeProjectInvitationUnprocessable)
}

func IsUnauthorizedError(err error) bool {
	return IsErrorWithCode(err, errCodeUnauthorizedInvalidToken) ||
		IsErrorWithCode(err, errCodeGithubIntegrationUnauthorized)
}

func IsForbiddenError(err error) bool {
	return IsErrorWithCode(err, errCodeForbiddenInsufficientUserRole) ||
		IsErrorWithCode(err, errCodeGithubIntegrationForbidden) ||
		IsErrorWithCode(err, errCodeForbiddenInsufficientProjectRole)
}

func IsConflictError(err error) bool {
	return IsErrorWithCode(err, errCodeEnvironmentDuplicateName) ||
		IsErrorWithCode(err, errCodeProjectInvitationAlreadyExists) ||
		IsErrorWithCode(err, errCodeProjectMemberAlreadyExists)
}
