package service

import (
	"context"
	"errors"

	reperr "release-manager/repository/errors"
	svcerr "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectMembershipManagementService struct {
	userSvc       model.UserService
	memberSvc     model.ProjectMemberService
	invitationSvc model.ProjectInvitationService
}

func (s *ProjectMembershipManagementService) Create(ctx context.Context, r model.ProjectMembershipRequest, requestedBy model.ProjectMember) (model.ProjectMembershipResponse, error) {

	if !requestedBy.CanGrantRole(r.Role) {
		return model.ProjectMembershipResponse{}, svcerr.ErrProjectMemberRoleCannotBeGranted
	}

	user, err := s.userSvc.GetByEmail(ctx, r.Email)

	// User already exists, add them to project
	if err == nil {
		return s.createMember(ctx, user.ID, r)
	}

	// User does not exist yet, send invitation
	if errors.Is(err, reperr.ErrResourceNotFound) {
		return s.createInvitation(ctx, r)
	}

	return model.ProjectMembershipResponse{}, err
}

func (s *ProjectMembershipManagementService) createInvitation(ctx context.Context, r model.ProjectMembershipRequest) (model.ProjectMembershipResponse, error) {
	i, err := s.invitationSvc.Create(ctx, r.ProjectID, r.Email, r.Role, r.RequestedByUserID)
	if err != nil {
		return model.ProjectMembershipResponse{}, err
	}

	return model.ProjectMembershipResponse{
		Status:   model.InvitationSentStatus,
		Resource: i,
	}, nil
}

func (s *ProjectMembershipManagementService) createMember(ctx context.Context, userID uuid.UUID, r model.ProjectMembershipRequest) (model.ProjectMembershipResponse, error) {
	i, err := s.memberSvc.Create(ctx, r.ProjectID, userID, r.Role, r.RequestedByUserID)
	if err != nil {
		return model.ProjectMembershipResponse{}, err
	}

	return model.ProjectMembershipResponse{
		Status:   model.MemberCreatedStatus,
		Resource: i,
	}, nil
}
