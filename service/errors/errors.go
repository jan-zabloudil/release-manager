package errors

import (
	"errors"
	"fmt"
)

var (
	ErrCodeUnauthenticatedUser             = "ERR_UNAUTHENTICATED_USER"
	ErrCodeInsufficientUserRole            = "ERR_INSUFFICIENT_USER_ROLE"
	ErrCodeInsufficientProjectRole         = "ERR_INSUFFICIENT_PROJECT_ROLE"
	ErrCodeUserNotProjectMember            = "ERR_USER_NOT_PROJECT_MEMBER"
	ErrCodeUserNotFound                    = "ERR_USER_NOT_FOUND"
	ErrCodeProjectNotFound                 = "ERR_PROJECT_NOT_FOUND"
	ErrCodeEnvironmentNotFound             = "ERR_ENVIRONMENT_NOT_FOUND"
	ErrCodeProjectInvalid                  = "ERR_PROJECT_INVALID"
	ErrCodeEnvironmentInvalid              = "ERR_ENVIRONMENT_INVALID"
	ErrCodeEnvironmentDuplicateName        = "ERR_ENVIRONMENT_DUPLICATE_NAME"
	ErrCodeSettingsUnprocessable           = "ERR_SETTINGS_UNPROCESSABLE"
	ErrCodeProjectInvitationInvalid        = "ERR_PROJECT_INVITATION_INVALID"
	ErrCodeProjectInvitationAlreadyExists  = "ERR_PROJECT_INVITATION_ALREADY_EXISTS"
	ErrCodeProjectInvitationNotFound       = "ERR_PROJECT_INVITATION_NOT_FOUND"
	ErrCodeProjectMemberAlreadyExists      = "ERR_PROJECT_MEMBER_ALREADY_EXISTS"
	ErrCodeGithubIntegrationNotEnabled     = "ERR_GITHUB_INTEGRATION_NOT_ENABLED"
	ErrCodeGithubClientUnauthorized        = "ERR_GITHUB_CLIENT_UNAUTHORIZED"
	ErrCodeGithubClientForbidden           = "ERR_GITHUB_CLIENT_FORBIDDEN"
	ErrCodeGithubRepoNotSetForProject      = "ERR_GITHUB_REPO_NOT_SET_FOR_PROJECT"
	ErrCodeGithubRepoNotFound              = "ERR_GITHUB_REPO_NOT_FOUND"
	ErrCodeGithubRepoInvalidURL            = "ERR_GITHUB_REPO_INVALID_URL"
	ErrCodeProjectMemberNotFound           = "ERR_PROJECT_MEMBER_NOT_FOUND"
	ErrCodeProjectMemberUnprocessable      = "ERR_PROJECT_MEMBER_UNPROCESSABLE"
	ErrCodeReleaseInvalid                  = "ERR_RELEASE_INVALID"
	ErrCodeReleaseNotFound                 = "ERR_RELEASE_NOT_FOUND"
	ErrCodeSlackIntegrationNotEnabled      = "ERR_SLACK_INTEGRATION_NOT_ENABLED"
	ErrCodeSlackClientUnauthorized         = "ERR_SLACK_CLIENT_UNAUTHORIZED"
	ErrCodeSlackChannelNotFound            = "ERR_SLACK_CHANNEL_NOT_FOUND"
	ErrCodeSlackChannelNotSetForProject    = "ERR_SLACK_CHANNEL_NOT_SET_FOR_PROJECT"
	ErrCodeGitTagNotFound                  = "ERR_GIT_TAG_NOT_FOUND"
	ErrCodeGithubReleaseNotFound           = "ERR_GITHUB_RELEASE_NOT_FOUND"
	ErrCodeReleaseGitTagAlreadyUsed        = "ERR_RELEASE_GIT_TAG_ALREADY_USED"
	ErrCodeDeploymentInvalid               = "ERR_DEPLOYMENT_INVALID"
	ErrCodeDeploymentNotFound              = "ERR_DEPLOYMENT_NOT_FOUND"
	ErrCodeProjectGithubRepoAlreadyUsed    = "ERR_PROJECT_GITHUB_REPO_ALREADY_USED"
	ErrCodeGithubNotesInvalidInput         = "ERR_GITHUB_NOTES_INVALID_INPUT"
	ErrCodeAdminUserCannotBeDeleted        = "ERR_ADMIN_USER_CANNOT_BE_DELETED"
	ErrCodeInvalidGithubTagDeletionWebhook = "ERR_INVALID_GITHUB_TAG_DELETION_WEBHOOK"
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

func NewProjectInvalidError() *Error {
	return &Error{
		Code:    ErrCodeProjectInvalid,
		Message: "Invalid project",
	}
}

func NewEnvironmentInvalidError() *Error {
	return &Error{
		Code:    ErrCodeEnvironmentInvalid,
		Message: "Invalid environment",
	}
}

func NewSettingsUnprocessableError() *Error {
	return &Error{
		Code:    ErrCodeSettingsUnprocessable,
		Message: "Settings unprocessable",
	}
}

func NewUnauthenticatedUserError() *Error {
	return &Error{
		Code:    ErrCodeUnauthenticatedUser,
		Message: "Unauthenticated user, user not found.",
	}
}

func NewInsufficientUserRoleError() *Error {
	return &Error{
		Code:    ErrCodeInsufficientUserRole,
		Message: "Forbidden access, insufficient user role.",
	}
}

func NewInsufficientProjectRoleError() *Error {
	return &Error{
		Code:    ErrCodeInsufficientProjectRole,
		Message: "Forbidden access, insufficient project role.",
	}
}

func NewUserNotProjectMemberError() *Error {
	return &Error{
		Code:    ErrCodeUserNotProjectMember,
		Message: "User is not a project member.",
	}
}

func NewProjectInvitationInvalidError() *Error {
	return &Error{
		Code:    ErrCodeProjectInvitationInvalid,
		Message: "Invalid project invitation",
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

func NewGithubRepoInvalidURL() *Error {
	return &Error{
		Code:    ErrCodeGithubRepoInvalidURL,
		Message: "Invalid Github repo URL.",
	}
}

func NewGithubIntegrationNotEnabledError() *Error {
	return &Error{
		Code:    ErrCodeGithubIntegrationNotEnabled,
		Message: "Github integration is not enabled.",
	}
}

func NewGithubRepoNotSetForProjectError() *Error {
	return &Error{
		Code:    ErrCodeGithubRepoNotSetForProject,
		Message: "Github repo is not set for the project.",
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

func NewGithubRepoNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeGithubRepoNotFound,
		Message: "Github repo not found among accessible repos.",
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

func NewReleaseInvalidError() *Error {
	return &Error{
		Code:    ErrCodeReleaseInvalid,
		Message: "Invalid release",
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

func NewSlackClientUnauthorizedError() *Error {
	return &Error{
		Code:    ErrCodeSlackClientUnauthorized,
		Message: "Cannot send Slack message, client is not properly authenticated (invalid or expired token).",
	}
}

func NewSlackChannelNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeSlackChannelNotFound,
		Message: "Slack channel not found.",
	}
}

func NewSlackChannelNotSetForProjectError() *Error {
	return &Error{
		Code:    ErrCodeSlackChannelNotSetForProject,
		Message: "Slack channel is not set for the project.",
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

func NewReleaseGitTagAlreadyUsedError() *Error {
	return &Error{
		Code:    ErrCodeReleaseGitTagAlreadyUsed,
		Message: "Git tag is already used for another release",
	}
}

func NewDeploymentInvalidError() *Error {
	return &Error{
		Code:    ErrCodeDeploymentInvalid,
		Message: "Invalid deployment",
	}
}

func NewDeploymentNotFoundError() *Error {
	return &Error{
		Code:    ErrCodeDeploymentNotFound,
		Message: "Deployment not found",
	}
}

func NewProjectGithubRepoAlreadyUsedError() *Error {
	return &Error{
		Code:    ErrCodeProjectGithubRepoAlreadyUsed,
		Message: "Github repo is already used for another project.",
	}
}

func NewGithubNotesInvalidInputError() *Error {
	return &Error{
		Code:    ErrCodeGithubNotesInvalidInput,
		Message: "Invalid input for generating release notes",
	}
}

func NewAdminUserCannotBeDeletedError() *Error {
	return &Error{
		Code:    ErrCodeAdminUserCannotBeDeleted,
		Message: "Admin user cannot be deleted",
	}
}

func NewInvalidGithubTagDeletionWebhookError() *Error {
	return &Error{
		Code:    ErrCodeInvalidGithubTagDeletionWebhook,
		Message: "Invalid Github webhook for tag deleted event",
	}
}

func IsErrorWithCode(err error, code string) bool {
	var svcErr *Error
	if errors.As(err, &svcErr) {
		return svcErr.Code == code
	}

	return false
}
