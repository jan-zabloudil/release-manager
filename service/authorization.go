package service

import (
	"context"
	"fmt"

	"release-manager/pkg/id"
	svcerrors "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type AuthorizationService struct {
	userRepo    userRepository
	projectRepo projectRepository
	releaseRepo releaseRepository
}

func NewAuthorizationService(userRepo userRepository, projectRepo projectRepository, releaseRepo releaseRepository) *AuthorizationService {
	return &AuthorizationService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
		releaseRepo: releaseRepo,
	}
}

func (s *AuthorizationService) AuthorizeUserRoleUser(ctx context.Context, userID id.AuthUser) error {
	return s.authorizeUserRole(ctx, userID, model.UserRoleUser)
}

func (s *AuthorizationService) AuthorizeUserRoleAdmin(ctx context.Context, userID id.AuthUser) error {
	return s.authorizeUserRole(ctx, userID, model.UserRoleAdmin)
}

func (s *AuthorizationService) AuthorizeProjectRoleEditor(ctx context.Context, projectID uuid.UUID, userID id.AuthUser) error {
	return s.authorizeProjectRole(ctx, projectID, userID, model.ProjectRoleEditor)
}

func (s *AuthorizationService) AuthorizeProjectRoleViewer(ctx context.Context, projectID uuid.UUID, userID id.AuthUser) error {
	return s.authorizeProjectRole(ctx, projectID, userID, model.ProjectRoleViewer)
}

// AuthorizeReleaseEditor authorizes the user that has editor role or higher in the release's project.
// the function can be used to authorize release action in the release service if the project ID is not directly provided.
func (s *AuthorizationService) AuthorizeReleaseEditor(ctx context.Context, releaseID uuid.UUID, userID id.AuthUser) error {
	return s.authorizeProjectRoleByRelease(ctx, releaseID, userID, model.ProjectRoleEditor)
}

// AuthorizeReleaseViewer authorizes the user that has viewer role or higher in the release's project.
// the function can be used to authorize release action in the release service if the project ID is not directly provided.
func (s *AuthorizationService) AuthorizeReleaseViewer(ctx context.Context, releaseID uuid.UUID, userID id.AuthUser) error {
	return s.authorizeProjectRoleByRelease(ctx, releaseID, userID, model.ProjectRoleViewer)
}

// authorizeProjectRoleByRelease checks if the user has the required or higher role in the project of the release.
func (s *AuthorizationService) authorizeProjectRoleByRelease(ctx context.Context, releaseID uuid.UUID, userID id.AuthUser, role model.ProjectRole) error {
	// Approach of reading the release and then calling authorizeProjectRole was chosen rather than
	// having repo function ReadProjectMemberByReleaseID, because:
	//
	// 1. Project repo should not access release.
	// 2. Reading release in separate query is very cheap operation.
	// 3. authorizeProjectRole function can be reused.
	rls, err := s.releaseRepo.ReadRelease(ctx, releaseID)
	if err != nil {
		return fmt.Errorf("reading release: %w", err)
	}

	return s.authorizeProjectRole(ctx, rls.ProjectID, userID, role)
}

// authorizeProjectRole checks if the user has the required or higher role in the project.
// If the user is not a member of the project, it checks if the user has admin role.
// If user is not a member with required role (or higher) and not an admin, it returns an error (ErrCodeUserNotProjectMember or ErrCodeInsufficientProjectRole).
// If project does not exist, it returns an error (ErrCodeProjectNotFound).
func (s *AuthorizationService) authorizeProjectRole(ctx context.Context, projectID uuid.UUID, userID id.AuthUser, role model.ProjectRole) error {
	member, err := s.projectRepo.ReadMember(ctx, projectID, uuid.UUID(userID))
	if err != nil {
		switch {
		case svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectMemberNotFound):
			// User is not a member of the project

			// First check if the project exists
			// Important to check if the project exists first otherwise we would authorize admin user even for non-existing projects
			if _, err := s.projectRepo.ReadProject(ctx, projectID); err != nil {
				return fmt.Errorf("reading project: %w", err)
			}

			// if project exists, check if the user is an admin (admin has access to all projects)
			if user, err := s.getUser(ctx, userID); err != nil {
				return fmt.Errorf("checking if user has admin role: %w", err)
			} else if user.IsAdmin() {
				return nil
			}

			// Project exists but user is not a member and also not an admin
			return svcerrors.NewUserNotProjectMemberError().Wrap(err)
		default:
			return fmt.Errorf("reading project member: %w", err)
		}
	}

	if !member.SatisfiesRequiredRole(role) {
		return svcerrors.NewInsufficientProjectRoleError()
	}

	return nil
}

func (s *AuthorizationService) authorizeUserRole(ctx context.Context, userID id.AuthUser, role model.UserRole) error {
	u, err := s.getUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if !u.HasAtLeastRole(role) {
		return svcerrors.NewInsufficientUserRoleError()
	}

	return nil
}

func (s *AuthorizationService) getUser(ctx context.Context, userID id.AuthUser) (model.User, error) {
	// Cannot use GetAuthenticated function from UserService because it would create a circular dependency.
	u, err := s.userRepo.Read(ctx, uuid.UUID(userID))
	if err != nil {
		switch {
		case svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeUserNotFound):
			return model.User{}, svcerrors.NewUnauthenticatedUserError().Wrap(err)
		default:
			return model.User{}, fmt.Errorf("reading user: %w", err)
		}
	}

	return u, nil
}
