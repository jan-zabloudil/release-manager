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
	authRepo    authRepository
	userRepo    userRepository
	projectRepo projectRepository
}

func NewAuthService(authRepo authRepository, userRepo userRepository, projectRepo projectRepository) *AuthService {
	return &AuthService{
		authRepo:    authRepo,
		userRepo:    userRepo,
		projectRepo: projectRepo,
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

func (s *AuthService) AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error {
	return s.AuthorizeUserRole(ctx, userID, model.UserRoleAdmin)
}

func (s *AuthService) AuthorizeUserRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error {
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

func (s *AuthService) AuthorizeProjectRoleEditor(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	return s.AuthorizeProjectRole(ctx, projectID, userID, model.ProjectRoleEditor)
}

func (s *AuthService) AuthorizeProjectRoleViewer(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	return s.AuthorizeProjectRole(ctx, projectID, userID, model.ProjectRoleViewer)
}

func (s *AuthService) AuthorizeProjectRole(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, role model.ProjectRole) error {
	user, err := s.userRepo.Read(ctx, userID)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return apierrors.NewUnauthorizedError().Wrap(err)
		default:
			return err
		}
	}

	// Admin user can access all projects
	if user.IsAdmin() {
		return nil
	}

	member, err := s.projectRepo.ReadMember(ctx, projectID, userID)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return apierrors.NewForbiddenInsufficientProjectRoleError().Wrap(err)
		default:
			return err
		}
	}

	if !member.HasAtLeastProjectRole(role) {
		return apierrors.NewForbiddenInsufficientProjectRoleError()
	}

	return nil
}
