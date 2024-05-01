package service

import (
	"context"

	"release-manager/pkg/apierrors"
	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/dberrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectMembershipService struct {
	authGuard        authGuard
	projectGetter    projectGetter
	repository       projectInvitationRepository
	invitationSender projectInvitationSender
}

func NewProjectMembershipService(
	guard authGuard,
	projectGetter projectGetter,
	repo projectInvitationRepository,
	invitationSender projectInvitationSender,
) *ProjectMembershipService {
	return &ProjectMembershipService{
		authGuard:        guard,
		projectGetter:    projectGetter,
		repository:       repo,
		invitationSender: invitationSender,
	}
}

func (s *ProjectMembershipService) CreateInvitation(
	ctx context.Context,
	c model.CreateProjectInvitationInput,
	authUserID uuid.UUID,
) (model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return model.ProjectInvitation{}, err
	}

	p, err := s.projectGetter.Get(ctx, c.ProjectID, authUserID)
	if err != nil {
		return model.ProjectInvitation{}, err
	}

	tkn, err := cryptox.NewToken()
	if err != nil {
		return model.ProjectInvitation{}, err
	}

	i, err := model.NewProjectInvitation(c, tkn, authUserID)
	if err != nil {
		return model.ProjectInvitation{}, apierrors.NewProjectInvitationUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	// TODO check if the user is already a member of the project

	if exists, err := s.invitationExists(ctx, i.Email, c.ProjectID); err != nil {
		return model.ProjectInvitation{}, err
	} else if exists {
		return model.ProjectInvitation{}, apierrors.NewProjectInvitationAlreadyExistsError()
	}

	if err := s.repository.Create(ctx, i); err != nil {
		return model.ProjectInvitation{}, err
	}

	s.invitationSender.SendProjectInvitation(ctx, model.ProjectInvitationInput{
		ProjectName:    p.Name,
		RecipientEmail: i.Email,
		Token:          tkn,
	})

	return i, nil
}

func (s *ProjectMembershipService) ListInvitations(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return nil, err
	}

	_, err := s.projectGetter.Get(ctx, projectID, authUserID)
	if err != nil {
		return nil, err
	}

	return s.repository.ReadAllForProject(ctx, projectID)
}

func (s *ProjectMembershipService) DeleteInvitation(ctx context.Context, projectID, invitationID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return err
	}

	_, err := s.projectGetter.Get(ctx, projectID, authUserID)
	if err != nil {
		return err
	}

	i, err := s.repository.Read(ctx, invitationID)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return apierrors.NewProjectInvitationNotFoundError().Wrap(err)
		default:
			return err
		}
	}

	return s.repository.Delete(ctx, i.ID)
}

func (s *ProjectMembershipService) AcceptInvitation(ctx context.Context, tkn cryptox.Token) error {
	invitation, err := s.getPendingInvitationByToken(ctx, tkn)
	if err != nil {
		return err
	}

	/*
		When an invitation is confirmed, db function handle_confirmed_invitation is triggered.
		The function checks if the email is already associated with a user account.
		If user exists, the function creates a project member and deletes the invitation.

		When user signs up, function handle_new_user is triggered.
		The function checks if the user's email is associated with a confirmed invitation(s).
		If so, the function creates a project member(s) and deletes the invitation(s).

		This approach is used because current repository implementation does not support transactions.
	*/
	invitation.Accept()
	return s.repository.Update(ctx, invitation)
}

func (s *ProjectMembershipService) RejectInvitation(ctx context.Context, tkn cryptox.Token) error {
	invitation, err := s.getPendingInvitationByToken(ctx, tkn)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, invitation.ID)
}

func (s *ProjectMembershipService) getPendingInvitationByToken(ctx context.Context, tkn cryptox.Token) (model.ProjectInvitation, error) {
	i, err := s.repository.ReadByTokenHashAndStatus(ctx, tkn.ToHash(), model.InvitationStatusPending)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return model.ProjectInvitation{}, apierrors.NewProjectInvitationNotFoundError().Wrap(err)
		default:
			return model.ProjectInvitation{}, err
		}
	}

	return i, nil
}

func (s *ProjectMembershipService) invitationExists(ctx context.Context, email string, projectID uuid.UUID) (bool, error) {
	if _, err := s.repository.ReadByEmailForProject(ctx, email, projectID); err != nil {
		if dberrors.IsNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
