package service

import (
	"context"

	"release-manager/pkg/apierrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type AuthorizationService struct {
	userRepo    userRepository
	projectRepo projectRepository
}

func NewAuthorizationService(userRepo userRepository, projectRepo projectRepository) *AuthorizationService {
	return &AuthorizationService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

func (s *AuthorizationService) AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error {
	return s.AuthorizeUserRole(ctx, userID, model.UserRoleAdmin)
}

func (s *AuthorizationService) AuthorizeUserRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error {
	user, err := s.userRepo.Read(ctx, userID)
	if err != nil {
		switch {
		case apierrors.IsNotFoundError(err):
			return apierrors.NewUnauthorizedUnknownUserError().Wrap(err)
		default:
			return err
		}
	}

	if !user.HasAtLeastRole(role) {
		return apierrors.NewForbiddenInsufficientUserRoleError()
	}

	return nil
}

func (s *AuthorizationService) AuthorizeProjectRoleEditor(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	return s.AuthorizeProjectRole(ctx, projectID, userID, model.ProjectRoleEditor)
}

func (s *AuthorizationService) AuthorizeProjectRoleViewer(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	return s.AuthorizeProjectRole(ctx, projectID, userID, model.ProjectRoleViewer)
}

func (s *AuthorizationService) AuthorizeProjectRole(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, role model.ProjectRole) error {
	user, err := s.userRepo.Read(ctx, userID)
	if err != nil {
		switch {
		case apierrors.IsNotFoundError(err):
			return apierrors.NewUnauthorizedUnknownUserError().Wrap(err)
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
		case apierrors.IsNotFoundError(err):
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
