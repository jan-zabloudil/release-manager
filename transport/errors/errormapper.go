package errors

import (
	"errors"

	svcerrors "release-manager/service/errors"
)

func ToError(err error) *Error {
	switch {
	case IsUnauthorizedError(err):
		return NewUnauthorizedError().Wrap(err)
	case IsForbiddenError(err):
		return NewForbiddenError().Wrap(err)
	case IsNotFoundError(err):
		return NewNotFoundError().Wrap(err)
	case IsUnprocessableModelError(err):
		return NewUnprocessableEntityError().Wrap(err)
	case IsConflictError(err):
		return NewConflictError().Wrap(err)
	case IsBadRequestError(err):
		return NewBadRequestError().Wrap(err)
	default:
		return NewServerError().Wrap(err)
	}
}

func IsNotFoundError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeUserNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeEnvironmentNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectInvitationNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubRepoNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubRepoNotSetForProject) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubIntegrationNotEnabled) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectMemberNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubIntegrationNotEnabled) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubRepoInvalidURL) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeReleaseNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGitTagNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubReleaseNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeSlackChannelNotFound)
}

func IsUnprocessableModelError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeProjectUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeEnvironmentUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeSettingsUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectInvitationUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectMemberUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeReleaseUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeDeploymentUnprocessable)
}

func IsUnauthorizedError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeUnauthorizedUnknownUser) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubClientUnauthorized) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeSlackClientUnauthorized)
}

func IsForbiddenError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeInsufficientUserRole) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubClientForbidden) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeInsufficientProjectRole) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeUserNotProjectMember) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeAdminUserCannotBeDeleted)
}

func IsConflictError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeEnvironmentDuplicateName) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectInvitationAlreadyExists) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectMemberAlreadyExists) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeReleaseGitTagAlreadyUsed) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectGithubRepoAlreadyUsed)
}

func IsBadRequestError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeSlackChannelNotSetForProject) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeSlackIntegrationNotEnabled) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubGeneratedNotesInvalidInput)
}

func IsGithubIntegrationNotEnabledError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeGithubIntegrationNotEnabled)
}

func isSvcErrorWithCode(err error, code string) bool {
	var svcErr *svcerrors.Error
	if errors.As(err, &svcErr) {
		return svcErr.Code == code
	}

	return false
}
