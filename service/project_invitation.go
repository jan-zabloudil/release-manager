package service

import (
	"context"
	"errors"

	reperr "release-manager/repository/errors"
	svcerr "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectInvitationService struct {
	repository model.ProjectInvitationRepository
}

func (s *ProjectInvitationService) Create(ctx context.Context, projectID uuid.UUID, email string, role model.ProjectRole, invitedByUserID uuid.UUID) (model.ProjectInvitation, error) {
	_, err := s.repository.ReadByEmail(ctx, projectID, email)
	if err == nil {
		return model.ProjectInvitation{}, svcerr.ErrInvitationAlreadyExists
	}

	if !errors.Is(err, reperr.ErrResourceNotFound) {
		return model.ProjectInvitation{}, err
	}

	return s.repository.Insert(ctx, projectID, email, role, invitedByUserID)
}

func (s *ProjectInvitationService) GetByEmail(ctx context.Context, projectID uuid.UUID, email string) (model.ProjectInvitation, error) {
	return s.repository.ReadByEmail(ctx, projectID, email)
}

func (s *ProjectInvitationService) Get(ctx context.Context, invitationID uuid.UUID) (model.ProjectInvitation, error) {
	return s.repository.Read(ctx, invitationID)
}

func (s *ProjectInvitationService) ListAll(ctx context.Context, projectID uuid.UUID) ([]model.ProjectInvitation, error) {
	return s.repository.ReadAll(ctx, projectID)
}

func (s *ProjectInvitationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
