package errors

import (
	"errors"

	"release-manager/pkg/validatorx"
	svcerrors "release-manager/service/errors"
)

// NewFromBodyUnmarshalErr creates an error from an error that occurred during unmarshalling the request body.
func NewFromBodyUnmarshalErr(err error) *Error {
	var validationErr validatorx.ValidationErrors
	if errors.As(err, &validationErr) {
		return NewInvalidRequestPayloadError().Wrap(err).WithData(validationErr).WithMessage("Validation errors")
	}

	return NewInvalidRequestPayloadError().Wrap(err).WithMessage(err.Error())
}

func NewFromSvcErr(err error) *Error {
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
	case isBadRequestError(err):
		return NewBadRequestError().Wrap(err)
	default:
		return NewServerError().Wrap(err)
	}
}

func isNotFoundError(err error) bool {
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

func isUnprocessableModelError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeEnvironmentUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSettingsUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectInvitationUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectMemberUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeReleaseUnprocessable) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeDeploymentUnprocessable)
}

func isUnauthorizedError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeUnauthenticatedUser) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubClientUnauthorized) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSlackClientUnauthorized)
}

func isForbiddenError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeInsufficientUserRole) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubClientForbidden) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeInsufficientProjectRole) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeUserNotProjectMember) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeAdminUserCannotBeDeleted)
}

func isConflictError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeEnvironmentDuplicateName) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectInvitationAlreadyExists) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectMemberAlreadyExists) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeReleaseGitTagAlreadyUsed) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectGithubRepoAlreadyUsed)
}

func isBadRequestError(err error) bool {
	return svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSlackChannelNotSetForProject) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeSlackIntegrationNotEnabled) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeGithubNotesInvalidInput) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeInvalidGithubTagDeletionWebhook) ||
		svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeInvalidProject)
}
