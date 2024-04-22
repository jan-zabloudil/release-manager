package util

import (
	"release-manager/pkg/apierrors"
	"release-manager/pkg/responseerrors"
)

func ToResponseError(err error) *responseerrors.ResponseError {
	switch {
	case apierrors.IsUnauthorizedError(err):
		return responseerrors.NewUnauthorizedError().Wrap(err)
	case apierrors.IsForbiddenError(err):
		return responseerrors.NewForbiddenError().Wrap(err)
	case apierrors.IsNotFoundError(err):
		return responseerrors.NewNotFoundError().Wrap(err)
	case apierrors.IsUnprocessableModelError(err):
		return responseerrors.NewUnprocessableEntityError().Wrap(err)
	case apierrors.IsConflictError(err):
		return responseerrors.NewConflictError().Wrap(err)
	default:
		return responseerrors.NewServerError().Wrap(err)
	}
}
