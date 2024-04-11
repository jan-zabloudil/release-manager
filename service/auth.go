package service

import (
	"context"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/autherrors"
	"release-manager/pkg/dberrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepo model.AuthRepository
	userRepo model.UserRepository
}

func NewAuthService(authRepo model.AuthRepository, userRepo model.UserRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	id, err := s.authRepo.ReadUserIDForToken(ctx, token)
	if err != nil {
		switch {
		case autherrors.IsInvalidTokenError(err):
			return uuid.Nil, apierrors.NewUnauthorizedError().Wrap(err)
		default:
			return uuid.Nil, err
		}
	}

	return id, nil
}

func (s *AuthService) AuthorizeAdminRole(ctx context.Context, userID uuid.UUID) error {
	return s.AuthorizeRole(ctx, userID, model.UserRoleAdmin)
}

func (s *AuthService) AuthorizeRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error {
	user, err := s.userRepo.Read(ctx, userID)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return apierrors.NewUnauthorizedError().Wrap(err)
		default:
			return err
		}
	}

	if !user.HasAtLeastRole(role) {
		return apierrors.NewForbiddenInsufficientUserRoleError()
	}

	return nil
}
