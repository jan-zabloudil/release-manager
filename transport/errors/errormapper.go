package errors

import (
	"errors"

	"release-manager/pkg/apierrors"
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
	return isAPIErrorWithCode(err, apierrors.ErrCodeUserNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeProjectNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeEnvironmentNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeProjectInvitationNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubRepositoryNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubRepositoryNotConfiguredForProject) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubIntegrationNotEnabled) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeProjectMemberNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubIntegrationNotEnabled) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubRepositoryInvalidURL) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeReleaseNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGitTagNotFound) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubReleaseNotFound)
}

func isUnprocessableModelError(err error) bool {
	return isAPIErrorWithCode(err, apierrors.ErrCodeProjectUnprocessable) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeEnvironmentUnprocessable) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeSettingsUnprocessable) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeProjectInvitationUnprocessable) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeProjectMemberUnprocessable) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeReleaseUnprocessable)
}

func isUnauthorizedError(err error) bool {
	return isAPIErrorWithCode(err, apierrors.ErrCodeUnauthorizedUnknownUser) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubClientUnauthorized)
}

func isForbiddenError(err error) bool {
	return isAPIErrorWithCode(err, apierrors.ErrCodeForbiddenInsufficientUserRole) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeGithubClientForbidden) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeForbiddenInsufficientProjectRole) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeForbiddenUserNotProjectMember)
}

func isConflictError(err error) bool {
	return isAPIErrorWithCode(err, apierrors.ErrCodeEnvironmentDuplicateName) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeProjectInvitationAlreadyExists) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeProjectMemberAlreadyExists) ||
		isAPIErrorWithCode(err, apierrors.ErrCodeReleaseDuplicateTitle)
}

func isAPIErrorWithCode(err error, code string) bool {
	var apiErr *apierrors.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == code
	}

	return false
}
