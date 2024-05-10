package apierrors

import (
	"errors"
	"fmt"
)

var (
	errCodeUnauthorizedUnknownUser                 = "ERR_UNAUTHORIZED_ACCESS_UNKNOWN_USER"
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
	errCodeGithubClientUnauthorized                = "ERR_GITHUB_CLIENT_UNAUTHORIZED"
	errCodeGithubClientForbidden                   = "ERR_GITHUB_CLIENT_FORBIDDEN"
	errCodeGithubRepositoryNotConfiguredForProject = "ERR_GITHUB_REPOSITORY_NOT_CONFIGURED_FOR_PROJECT"
	errCodeGithubRepositoryNotFound                = "ERR_GITHUB_REPOSITORY_NOT_FOUND"
	errCodeGithubRepositoryInvalidURL              = "ERR_GITHUB_REPOSITORY_INVALID_URL"
	errCodeProjectMemberNotFound                   = "ERR_PROJECT_MEMBER_NOT_FOUND"
	errCodeProjectMemberUnprocessable              = "ERR_PROJECT_MEMBER_UNPROCESSABLE"
	errCodeReleaseUnprocessable                    = "ERR_RELEASE_UNPROCESSABLE"
	errCodeReleaseNotFound                         = "ERR_RELEASE_NOT_FOUND"
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

func NewUnauthorizedUnknownUserError() *APIError {
	return &APIError{
		Code:    errCodeUnauthorizedUnknownUser,
		Message: "Unauthorized access, unknown user.",
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

func NewGithubRepositoryInvalidURL() *APIError {
	return &APIError{
		Code:    errCodeGithubRepositoryInvalidURL,
		Message: "Invalid Github repository URL.",
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

func NewGithubClientUnauthorizedError() *APIError {
	return &APIError{
		Code:    errCodeGithubClientUnauthorized,
		Message: "Request to the GitHub API cannot be processed because the client is not properly authenticated (invalid or expired token).",
	}
}

func NewGithubClientForbiddenError() *APIError {
	return &APIError{
		Code:    errCodeGithubClientForbidden,
		Message: "Request cannot be processed because the client does not have permission to access the specified resource via GitHub API.",
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

func NewProjectMemberUnprocessableError() *APIError {
	return &APIError{
		Code:    errCodeProjectMemberUnprocessable,
		Message: "Project member unprocessable",
	}
}

func NewReleaseUnprocessableError() *APIError {
	return &APIError{
		Code:    errCodeReleaseUnprocessable,
		Message: "Release unprocessable",
	}
}

func NewReleaseNotFoundError() *APIError {
	return &APIError{
		Code:    errCodeReleaseNotFound,
		Message: "Release not found",
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
		IsErrorWithCode(err, errCodeProjectMemberNotFound) ||
		IsErrorWithCode(err, errCodeGithubIntegrationNotEnabled) ||
		IsErrorWithCode(err, errCodeGithubRepositoryInvalidURL) ||
		IsErrorWithCode(err, errCodeReleaseNotFound)
}

func IsUnprocessableModelError(err error) bool {
	return IsErrorWithCode(err, errCodeProjectUnprocessable) ||
		IsErrorWithCode(err, errCodeEnvironmentUnprocessable) ||
		IsErrorWithCode(err, errCodeSettingsUnprocessable) ||
		IsErrorWithCode(err, errCodeProjectInvitationUnprocessable) ||
		IsErrorWithCode(err, errCodeProjectMemberUnprocessable) ||
		IsErrorWithCode(err, errCodeReleaseUnprocessable)
}

func IsUnauthorizedError(err error) bool {
	return IsErrorWithCode(err, errCodeUnauthorizedUnknownUser) ||
		IsErrorWithCode(err, errCodeGithubClientUnauthorized)
}

func IsForbiddenError(err error) bool {
	return IsErrorWithCode(err, errCodeForbiddenInsufficientUserRole) ||
		IsErrorWithCode(err, errCodeGithubClientForbidden) ||
		IsErrorWithCode(err, errCodeForbiddenInsufficientProjectRole)
}

func IsConflictError(err error) bool {
	return IsErrorWithCode(err, errCodeEnvironmentDuplicateName) ||
		IsErrorWithCode(err, errCodeProjectInvitationAlreadyExists) ||
		IsErrorWithCode(err, errCodeProjectMemberAlreadyExists)
}
