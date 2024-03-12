package service

import (
	"context"
	"errors"

	reperr "release-manager/repository/errors"
	svcerr "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectMemberService struct {
	repository model.ProjectMemberRepository
}

func (s *ProjectMemberService) Create(ctx context.Context, projectID, userID uuid.UUID, role model.ProjectRole, invitedByUserID uuid.UUID) (model.ProjectMember, error) {
	_, err := s.repository.Read(ctx, projectID, userID)
	if err == nil {
		return model.ProjectMember{}, svcerr.ErrUserIsAlreadyMember
	}

	if !errors.Is(err, reperr.ErrResourceNotFound) {
		return model.ProjectMember{}, err
	}

	return s.repository.Insert(ctx, projectID, userID, role, invitedByUserID)
}

func (s *ProjectMemberService) Get(ctx context.Context, projectID, userID uuid.UUID) (model.ProjectMember, error) {
	return s.repository.Read(ctx, projectID, userID)
}

func (s *ProjectMemberService) UpdateRole(ctx context.Context, member, updatedBy model.ProjectMember, newRole model.ProjectRole) (model.ProjectMember, error) {
	if !updatedBy.CanUpdateMember(member) {
		return model.ProjectMember{}, svcerr.ErrProjectMemberUpdateNotAllowed
	}

	if !updatedBy.CanGrantRole(newRole) {
		return model.ProjectMember{}, svcerr.ErrProjectMemberRoleCannotBeGranted
	}

	member.Role = newRole
	return s.repository.Update(ctx, member)
}

func (s *ProjectMemberService) Delete(ctx context.Context, projectID, userID uuid.UUID) error {
	return s.repository.Delete(ctx, projectID, userID)
}

func (s *ProjectMemberService) ListAll(ctx context.Context, projectID uuid.UUID) ([]model.ProjectMember, error) {
	return s.repository.ReadAll(ctx, projectID)
}
