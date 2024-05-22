package apierrors

import (
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrCodeUnauthorizedUnknownUser                 = "ERR_UNAUTHORIZED_ACCESS_UNKNOWN_USER"
	ErrCodeForbiddenInsufficientUserRole           = "ERR_FORBIDDEN_ACCESS_INSUFFICIENT_USER_ROLE"
	ErrCodeForbiddenInsufficientProjectRole        = "ERR_FORBIDDEN_ACCESS_INSUFFICIENT_PROJECT_ROLE"
	ErrCodeForbiddenUserNotProjectMember           = "ERR_FORBIDDEN_USER_NOT_PROJECT_MEMBER"
	ErrCodeUserNotFound                            = "ERR_USER_NOT_FOUND"
	ErrCodeProjectNotFound                         = "ERR_PROJECT_NOT_FOUND"
	ErrCodeEnvironmentNotFound                     = "ERR_ENVIRONMENT_NOT_FOUND"
	ErrCodeProjectUnprocessable                    = "ERR_PROJECT_UNPROCESSABLE"
	ErrCodeEnvironmentUnprocessable                = "ERR_ENVIRONMENT_UNPROCESSABLE"
	ErrCodeEnvironmentDuplicateName                = "ERR_ENVIRONMENT_DUPLICATE_NAME"
	ErrCodeSettingsUnprocessable                   = "ERR_SETTINGS_UNPROCESSABLE"
	ErrCodeProjectInvitationUnprocessable          = "ERR_PROJECT_INVITATION_UNPROCESSABLE"
	ErrCodeProjectInvitationAlreadyExists          = "ERR_PROJECT_INVITATION_ALREADY_EXISTS"
	ErrCodeProjectInvitationNotFound               = "ERR_PROJECT_INVITATION_NOT_FOUND"
	ErrCodeProjectMemberAlreadyExists              = "ERR_PROJECT_MEMBER_ALREADY_EXISTS"
	ErrCodeGithubIntegrationNotEnabled             = "ERR_GITHUB_INTEGRATION_NOT_ENABLED"
	ErrCodeGithubClientUnauthorized                = "ERR_GITHUB_CLIENT_UNAUTHORIZED"
	ErrCodeGithubClientForbidden                   = "ERR_GITHUB_CLIENT_FORBIDDEN"
	ErrCodeGithubRepositoryNotConfiguredForProject = "ERR_GITHUB_REPOSITORY_NOT_CONFIGURED_FOR_PROJECT"
	ErrCodeGithubRepositoryNotFound                = "ERR_GITHUB_REPOSITORY_NOT_FOUND"
	ErrCodeGithubRepositoryInvalidURL              = "ERR_GITHUB_REPOSITORY_INVALID_URL"
	ErrCodeProjectMemberNotFound                   = "ERR_PROJECT_MEMBER_NOT_FOUND"
	ErrCodeProjectMemberUnprocessable              = "ERR_PROJECT_MEMBER_UNPROCESSABLE"
	ErrCodeReleaseUnprocessable                    = "ERR_RELEASE_UNPROCESSABLE"
	ErrCodeReleaseNotFound                         = "ERR_RELEASE_NOT_FOUND"
	ErrCodeReleaseDuplicateTitle                   = "ERR_RELEASE_DUPLICATE_TITLE"
	ErrCodeSlackIntegrationNotEnabled              = "ERR_SLACK_INTEGRATION_NOT_ENABLED"
	ErrCodeGitTagNotFound                          = "ERR_GIT_TAG_NOT_FOUND"
	ErrCodeGithubReleaseNotFound                   = "ERR_GITHUB_RELEASE_NOT_FOUND"
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
		Code:    ErrCodeUserNotFound,
		Message: "User not found",
	}
}

func NewProjectNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectNotFound,
		Message: "Project not found",
	}
}

func NewEnvironmentNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeEnvironmentNotFound,
		Message: "Environment not found",
	}
}

func NewEnvironmentDuplicateNameError() *APIError {
	return &APIError{
		Code:    ErrCodeEnvironmentDuplicateName,
		Message: "environment name is already in use",
	}
}

func NewProjectUnprocessableError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectUnprocessable,
		Message: "Project unprocessable",
	}
}

func NewEnvironmentUnprocessableError() *APIError {
	return &APIError{
		Code:    ErrCodeEnvironmentUnprocessable,
		Message: "Environment unprocessable",
	}
}

func NewSettingsUnprocessableError() *APIError {
	return &APIError{
		Code:    ErrCodeSettingsUnprocessable,
		Message: "Settings unprocessable",
	}
}

func NewUnauthorizedUnknownUserError() *APIError {
	return &APIError{
		Code:    ErrCodeUnauthorizedUnknownUser,
		Message: "Unauthorized access, unknown user.",
	}
}

func NewForbiddenInsufficientUserRoleError() *APIError {
	return &APIError{
		Code:    ErrCodeForbiddenInsufficientUserRole,
		Message: "Forbidden access, insufficient user role.",
	}
}

func NewForbiddenInsufficientProjectRoleError() *APIError {
	return &APIError{
		Code:    ErrCodeForbiddenInsufficientProjectRole,
		Message: "Forbidden access, insufficient project role.",
	}
}

func NewForbiddenUserNotProjectMemberError() *APIError {
	return &APIError{
		Code:    ErrCodeForbiddenUserNotProjectMember,
		Message: "User is not a project member.",
	}
}

func NewProjectInvitationUnprocessableError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectInvitationUnprocessable,
		Message: "Project invitation unprocessable",
	}
}

func NewProjectInvitationAlreadyExistsError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectInvitationAlreadyExists,
		Message: "Project invitation already exists",
	}
}

func NewProjectInvitationNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectInvitationNotFound,
		Message: "Project invitation not found",
	}
}

func NewProjectMemberAlreadyExistsError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectMemberAlreadyExists,
		Message: "Project member already exists",
	}
}

func NewGithubRepositoryInvalidURL() *APIError {
	return &APIError{
		Code:    ErrCodeGithubRepositoryInvalidURL,
		Message: "Invalid Github repository URL.",
	}
}

func NewGithubIntegrationNotEnabledError() *APIError {
	return &APIError{
		Code:    ErrCodeGithubIntegrationNotEnabled,
		Message: "Github integration is not enabled.",
	}
}

func NewGithubRepositoryNotConfiguredForProjectError() *APIError {
	return &APIError{
		Code:    ErrCodeGithubRepositoryNotConfiguredForProject,
		Message: "Github repository is not configured for the project.",
	}
}

func NewGithubClientUnauthorizedError() *APIError {
	return &APIError{
		Code:    ErrCodeGithubClientUnauthorized,
		Message: "Request to the GitHub API cannot be processed because the client is not properly authenticated (invalid or expired token).",
	}
}

func NewGithubClientForbiddenError() *APIError {
	return &APIError{
		Code:    ErrCodeGithubClientForbidden,
		Message: "Request cannot be processed because the client does not have permission to access the specified resource via GitHub API.",
	}
}

func NewGithubRepositoryNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeGithubRepositoryNotFound,
		Message: "Github repository not found among accessible repositories.",
	}
}

func NewProjectMemberNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectMemberNotFound,
		Message: "Project member not found",
	}
}

func NewProjectMemberUnprocessableError() *APIError {
	return &APIError{
		Code:    ErrCodeProjectMemberUnprocessable,
		Message: "Project member unprocessable",
	}
}

func NewReleaseUnprocessableError() *APIError {
	return &APIError{
		Code:    ErrCodeReleaseUnprocessable,
		Message: "Release unprocessable",
	}
}

func NewReleaseNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeReleaseNotFound,
		Message: "Release not found",
	}
}

func NewReleaseDuplicateTitleError() *APIError {
	return &APIError{
		Code:    ErrCodeReleaseDuplicateTitle,
		Message: "Release with the same title already exists",
	}
}

func NewSlackIntegrationNotEnabledError() *APIError {
	return &APIError{
		Code:    ErrCodeSlackIntegrationNotEnabled,
		Message: "Slack integration is not enabled.",
	}
}

func NewGitTagNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeGitTagNotFound,
		Message: "Git tag not found",
	}
}

func NewGithubReleaseNotFoundError() *APIError {
	return &APIError{
		Code:    ErrCodeGithubReleaseNotFound,
		Message: "Github release not found",
	}
}

func IsNotFoundError(err error) bool {
	return isErrorWithCode(err, ErrCodeUserNotFound) ||
		isErrorWithCode(err, ErrCodeProjectNotFound) ||
		isErrorWithCode(err, ErrCodeEnvironmentNotFound) ||
		isErrorWithCode(err, ErrCodeProjectInvitationNotFound) ||
		isErrorWithCode(err, ErrCodeGithubRepositoryNotFound) ||
		isErrorWithCode(err, ErrCodeGithubRepositoryNotConfiguredForProject) ||
		isErrorWithCode(err, ErrCodeGithubIntegrationNotEnabled) ||
		isErrorWithCode(err, ErrCodeProjectMemberNotFound) ||
		isErrorWithCode(err, ErrCodeGithubIntegrationNotEnabled) ||
		isErrorWithCode(err, ErrCodeGithubRepositoryInvalidURL) ||
		isErrorWithCode(err, ErrCodeReleaseNotFound) ||
		isErrorWithCode(err, ErrCodeGitTagNotFound) ||
		isErrorWithCode(err, ErrCodeGithubReleaseNotFound)
}

func IsInsufficientUserRoleError(err error) bool {
	return isErrorWithCode(err, ErrCodeForbiddenInsufficientUserRole)
}

func IsSlackIntegrationNotEnabledError(err error) bool {
	return isErrorWithCode(err, ErrCodeSlackIntegrationNotEnabled)
}

func GetLogLevel(err error) slog.Level {
	if IsSlackIntegrationNotEnabledError(err) {
		return slog.LevelDebug
	}

	return slog.LevelError
}

func isErrorWithCode(err error, code string) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == code
	}

	return false
}
