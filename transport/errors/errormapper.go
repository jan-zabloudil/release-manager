package errors

import (
	"errors"

	svcerrors "release-manager/service/errors"
)

func ToError(err error) *Error {
	switch {
	case isUnauthorizedError(err):
		return NewUnauthorizedError().Wrap(err)
	case isForbiddenError(err):
		return NewForbiddenError().Wrap(err)
	case isNotFoundError(err):
		return NewNotFoundError().Wrap(err)
	case isUnprocessableModelError(err):
		return NewUnprocessableEntityError().Wrap(err)
	case isConflictError(err):
		return NewConflictError().Wrap(err)
	default:
		return NewServerError().Wrap(err)
	}
}

func isNotFoundError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeUserNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeEnvironmentNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectInvitationNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubRepositoryNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubRepositoryNotConfiguredForProject) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubIntegrationNotEnabled) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectMemberNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubIntegrationNotEnabled) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubRepositoryInvalidURL) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeReleaseNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGitTagNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubReleaseNotFound) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeSlackChannelNotFound)
}

func isUnprocessableModelError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeProjectUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeEnvironmentUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeSettingsUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectInvitationUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectMemberUnprocessable) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeReleaseUnprocessable)
}

func isUnauthorizedError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeUnauthorizedUnknownUser) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubClientUnauthorized) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeSlackClientUnauthorized)
}

func isForbiddenError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeForbiddenInsufficientUserRole) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeGithubClientForbidden) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeForbiddenInsufficientProjectRole) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeForbiddenUserNotProjectMember)
}

func isConflictError(err error) bool {
	return isSvcErrorWithCode(err, svcerrors.ErrCodeEnvironmentDuplicateName) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectInvitationAlreadyExists) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeProjectMemberAlreadyExists) ||
		isSvcErrorWithCode(err, svcerrors.ErrCodeReleaseDuplicateTitle)
}

func isSvcErrorWithCode(err error, code string) bool {
	var svcErr *svcerrors.Error
	if errors.As(err, &svcErr) {
		return svcErr.Code == code
	}

	return false
}
