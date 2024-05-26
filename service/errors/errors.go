package errors

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
	ErrCodeSlackIntegrationNotEnabled              = "ERR_SLACK_INTEGRATION_NOT_ENABLED"
	ErrCodeGitTagNotFound                          = "ERR_GIT_TAG_NOT_FOUND"
	ErrCodeGithubReleaseNotFound                   = "ERR_GITHUB_RELEASE_NOT_FOUND"
)

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %s, error: %s", e.Code, e.Err)
}

func (e *Error) Wrap(err error) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

func (e *Error) WithMessage(msg string) *Error {
	e.Message = msg
	return e
}

func NewUserNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeUserNotFound,
		Message: "User not found",
	}
}

func NewProjectNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeProjectNotFound,
		Message: "Project not found",
	}
}

func NewEnvironmentNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeEnvironmentNotFound,
		Message: "Environment not found",
	}
}

func NewEnvironmentDuplicateNameError() *Error {
	return &Error{
		Code:    ErrCodeEnvironmentDuplicateName,
		Message: "environment name is already in use",
	}
}

func NewProjectUnprocessableError() *Error {
	return &Error{
		Code:    ErrCodeProjectUnprocessable,
		Message: "Project unprocessable",
	}
}

func NewEnvironmentUnprocessableError() *Error {
	return &Error{
		Code:    ErrCodeEnvironmentUnprocessable,
		Message: "Environment unprocessable",
	}
}

func NewSettingsUnprocessableError() *Error {
	return &Error{
		Code:    ErrCodeSettingsUnprocessable,
		Message: "Settings unprocessable",
	}
}

func NewUnauthorizedUnknownUserError() *Error {
	return &Error{
		Code:    ErrCodeUnauthorizedUnknownUser,
		Message: "Unauthorized access, unknown user.",
	}
}

func NewForbiddenInsufficientUserRoleError() *Error {
	return &Error{
		Code:    ErrCodeForbiddenInsufficientUserRole,
		Message: "Forbidden access, insufficient user role.",
	}
}

func NewForbiddenInsufficientProjectRoleError() *Error {
	return &Error{
		Code:    ErrCodeForbiddenInsufficientProjectRole,
		Message: "Forbidden access, insufficient project role.",
	}
}

func NewForbiddenUserNotProjectMemberError() *Error {
	return &Error{
		Code:    ErrCodeForbiddenUserNotProjectMember,
		Message: "User is not a project member.",
	}
}

func NewProjectInvitationUnprocessableError() *Error {
	return &Error{
		Code:    ErrCodeProjectInvitationUnprocessable,
		Message: "Project invitation unprocessable",
	}
}

func NewProjectInvitationAlreadyExistsError() *Error {
	return &Error{
		Code:    ErrCodeProjectInvitationAlreadyExists,
		Message: "Project invitation already exists",
	}
}

func NewProjectInvitationNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeProjectInvitationNotFound,
		Message: "Project invitation not found",
	}
}

func NewProjectMemberAlreadyExistsError() *Error {
	return &Error{
		Code:    ErrCodeProjectMemberAlreadyExists,
		Message: "Project member already exists",
	}
}

func NewGithubRepositoryInvalidURL() *Error {
	return &Error{
		Code:    ErrCodeGithubRepositoryInvalidURL,
		Message: "Invalid Github repository URL.",
	}
}

func NewGithubIntegrationNotEnabledError() *Error {
	return &Error{
		Code:    ErrCodeGithubIntegrationNotEnabled,
		Message: "Github integration is not enabled.",
	}
}

func NewGithubRepositoryNotConfiguredForProjectError() *Error {
	return &Error{
		Code:    ErrCodeGithubRepositoryNotConfiguredForProject,
		Message: "Github repository is not configured for the project.",
	}
}

func NewGithubClientUnauthorizedError() *Error {
	return &Error{
		Code:    ErrCodeGithubClientUnauthorized,
		Message: "Request to the GitHub API cannot be processed because the client is not properly authenticated (invalid or expired token).",
	}
}

func NewGithubClientForbiddenError() *Error {
	return &Error{
		Code:    ErrCodeGithubClientForbidden,
		Message: "Request cannot be processed because the client does not have permission to access the specified resource via GitHub API.",
	}
}

func NewGithubRepositoryNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeGithubRepositoryNotFound,
		Message: "Github repository not found among accessible repositories.",
	}
}

func NewProjectMemberNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeProjectMemberNotFound,
		Message: "Project member not found",
	}
}

func NewProjectMemberUnprocessableError() *Error {
	return &Error{
		Code:    ErrCodeProjectMemberUnprocessable,
		Message: "Project member unprocessable",
	}
}

func NewReleaseUnprocessableError() *Error {
	return &Error{
		Code:    ErrCodeReleaseUnprocessable,
		Message: "Release unprocessable",
	}
}

func NewReleaseNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeReleaseNotFound,
		Message: "Release not found",
	}
}

func NewSlackIntegrationNotEnabledError() *Error {
	return &Error{
		Code:    ErrCodeSlackIntegrationNotEnabled,
		Message: "Slack integration is not enabled.",
	}
}

func NewGitTagNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeGitTagNotFound,
		Message: "Git tag not found",
	}
}

func NewGithubReleaseNotFoundError() *Error {
	return &Error{
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
	var svcErr *Error
	if errors.As(err, &svcErr) {
		return svcErr.Code == code
	}

	return false
}
