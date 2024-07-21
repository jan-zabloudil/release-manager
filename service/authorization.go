package service

import (
	"context"
	"fmt"

	svcerrors "release-manager/service/errors"
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

func (s *AuthorizationService) AuthorizeUserRoleUser(ctx context.Context, userID uuid.UUID) error {
	return s.authorizeUserRole(ctx, userID, model.UserRoleUser)
}

func (s *AuthorizationService) AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error {
	return s.authorizeUserRole(ctx, userID, model.UserRoleAdmin)
}

func (s *AuthorizationService) authorizeUserRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error {
	user, err := s.userRepo.Read(ctx, userID)
	if err != nil {
		switch {
		case svcerrors.IsNotFoundError(err):
			return svcerrors.NewUnauthorizedUnknownUserError().Wrap(err)
		default:
			return fmt.Errorf("reading user: %w", err)
		}
	}

	if !user.HasAtLeastRole(role) {
		return svcerrors.NewForbiddenInsufficientUserRoleError()
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
		case svcerrors.IsNotFoundError(err):
			return svcerrors.NewUnauthorizedUnknownUserError().Wrap(err)
		default:
			return fmt.Errorf("reading user: %w", err)
		}
	}

	// Admin user can access all projects
	if user.IsAdmin() {
		return nil
	}

	member, err := s.projectRepo.ReadMember(ctx, projectID, userID)
	if err != nil {
		switch {
		case svcerrors.IsNotFoundError(err):
			return svcerrors.NewForbiddenUserNotProjectMemberError().Wrap(err)
		default:
			return fmt.Errorf("reading project member: %w", err)
		}
	}

	if !member.HasAtLeastProjectRole(role) {
		return svcerrors.NewForbiddenInsufficientProjectRoleError()
	}

	return nil
}
