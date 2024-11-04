package errors

import (
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
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeUserNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeEnvironmentNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectInvitationNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubRepoNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubRepoNotSetForProject) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubIntegrationNotEnabled) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectMemberNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubIntegrationNotEnabled) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubRepoInvalidURL) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeReleaseNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGitTagNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubReleaseNotFound) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSlackChannelNotFound)
}

func IsUnprocessableModelError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeEnvironmentUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSettingsUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectInvitationUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectMemberUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeReleaseUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeDeploymentUnprocessable)
}

func IsUnauthorizedError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeUnauthenticatedUser) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubClientUnauthorized) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSlackClientUnauthorized)
}

func IsForbiddenError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeInsufficientUserRole) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubClientForbidden) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeInsufficientProjectRole) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeUserNotProjectMember) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeAdminUserCannotBeDeleted)
}

func IsConflictError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeEnvironmentDuplicateName) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectInvitationAlreadyExists) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectMemberAlreadyExists) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeReleaseGitTagAlreadyUsed) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectGithubRepoAlreadyUsed)
}

func IsBadRequestError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSlackChannelNotSetForProject) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSlackIntegrationNotEnabled) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubNotesInvalidInput)
}
