package service

import (
	"context"

	"release-manager/pkg/apierrors"
	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/dberrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectService struct {
	authGuard      authGuard
	settingsGetter settingsGetter
	userGetter     userGetter
	emailSender    emailSender
	githubManager  githubManager
	repo           projectRepository
}

func NewProjectService(
	guard authGuard,
	settingsGetter settingsGetter,
	userGetter userGetter,
	emailSender emailSender,
	githubManager githubManager,
	repo projectRepository,
) *ProjectService {
	return &ProjectService{
		authGuard:      guard,
		settingsGetter: settingsGetter,
		userGetter:     userGetter,
		emailSender:    emailSender,
		githubManager:  githubManager,
		repo:           repo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, c model.CreateProjectInput, authUserID uuid.UUID) (model.Project, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Project{}, err
	}

	p, err := model.NewProject(c)
	if err != nil {
		return model.Project{}, apierrors.NewProjectUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	u, err := s.userGetter.Get(ctx, authUserID, authUserID)
	if err != nil {
		return model.Project{}, err
	}

	owner, err := model.NewProjectOwner(u, p.ID)
	if err != nil {
		return model.Project{}, err
	}

	if err := s.repo.CreateProjectWithOwner(ctx, p, owner); err != nil {
		return model.Project{}, err
	}

	return p, nil
}

func (s *ProjectService) GetProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (model.Project, error) {
	// TODO add project member authorization

	p, err := s.repo.ReadProject(ctx, projectID)
	if err != nil {
		return model.Project{}, err
	}

	return p, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, authUserID uuid.UUID) ([]model.Project, error) {
	// TODO add project member authorization
	// TODO fetch only project where the user is a member

	return s.repo.ReadAllProjects(ctx)
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return err
	}

	err := s.repo.DeleteProject(ctx, projectID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, u model.UpdateProjectInput, projectID, authUserID uuid.UUID) (model.Project, error) {
	// TODO add project member authorization

	p, err := s.repo.UpdateProject(ctx, projectID, func(p model.Project) (model.Project, error) {
		if err := p.Update(u); err != nil {
			return model.Project{}, apierrors.NewProjectUnprocessableError().Wrap(err).WithMessage(err.Error())
		}

		return p, nil
	})
	if err != nil {
		return model.Project{}, err
	}

	return p, nil
}

func (s *ProjectService) CreateEnvironment(ctx context.Context, c model.CreateEnvironmentInput, authUserID uuid.UUID) (model.Environment, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Environment{}, err
	}

	_, err := s.GetProject(ctx, c.ProjectID, authUserID)
	if err != nil {
		return model.Environment{}, err
	}

	env, err := model.NewEnvironment(c)
	if err != nil {
		return model.Environment{}, apierrors.NewEnvironmentUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if isUnique, err := s.isEnvironmentNameUnique(ctx, env.ProjectID, env.Name); err != nil {
		return model.Environment{}, err
	} else if !isUnique {
		return model.Environment{}, apierrors.NewEnvironmentDuplicateNameError()
	}

	if err := s.repo.CreateEnvironment(ctx, env); err != nil {
		return model.Environment{}, err
	}

	return env, nil
}

func (s *ProjectService) GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (model.Environment, error) {
	// TODO add project member authorization

	env, err := s.repo.ReadEnvironment(ctx, projectID, envID)
	if err != nil {
		return model.Environment{}, err
	}

	return env, nil
}

func (s *ProjectService) UpdateEnvironment(ctx context.Context, u model.UpdateEnvironmentInput, projectID, envID, authUserID uuid.UUID) (model.Environment, error) {
	// TODO add project member authorization

	if _, err := s.GetProject(ctx, projectID, authUserID); err != nil {
		return model.Environment{}, err
	}

	env, err := s.GetEnvironment(ctx, projectID, envID, authUserID)
	if err != nil {
		return model.Environment{}, err
	}

	originalName := env.Name
	if err := env.Update(u); err != nil {
		return model.Environment{}, apierrors.NewEnvironmentUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if originalName != env.Name {
		if isUnique, err := s.isEnvironmentNameUnique(ctx, env.ProjectID, env.Name); err != nil {
			return model.Environment{}, err
		} else if !isUnique {
			return model.Environment{}, apierrors.NewEnvironmentDuplicateNameError()
		}
	}

	if err := s.repo.UpdateEnvironment(ctx, env); err != nil {
		return model.Environment{}, err
	}

	return env, nil
}

func (s *ProjectService) ListEnvironments(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.Environment, error) {
	// TODO add project member authorization

	envs, err := s.repo.ListEnvironmentsForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if len(envs) == 0 {
		exists, err := s.projectExists(ctx, projectID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apierrors.NewProjectNotFoundError()
		}
	}

	return envs, nil
}

func (s *ProjectService) DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return err
	}

	err := s.repo.DeleteEnvironment(ctx, projectID, envID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) ListGithubRepositoryTags(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.GitTag, error) {
	// TODO add project member authorization

	project, err := s.GetProject(ctx, projectID, authUserID)
	if err != nil {
		return nil, err
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return nil, err
	}

	if !project.IsGithubConfigured() {
		return nil, apierrors.NewGithubRepositoryNotConfiguredForProjectError()
	}

	return s.githubManager.ReadTagsForRepository(ctx, tkn, project.GithubRepositoryURL)
}

func (s *ProjectService) Invite(ctx context.Context, c model.CreateProjectInvitationInput, authUserID uuid.UUID) (model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.ProjectInvitation{}, err
	}

	p, err := s.GetProject(ctx, c.ProjectID, authUserID)
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

	memberExists, err := s.memberExists(ctx, i.ProjectID, i.Email)
	if err != nil {
		return model.ProjectInvitation{}, err
	}

	if memberExists {
		return model.ProjectInvitation{}, apierrors.NewProjectMemberAlreadyExistsError()
	}

	invitationExists, err := s.invitationExists(ctx, i.Email, c.ProjectID)
	if err != nil {
		return model.ProjectInvitation{}, err
	}
	if invitationExists {
		return model.ProjectInvitation{}, apierrors.NewProjectInvitationAlreadyExistsError()
	}

	if err := s.repo.CreateInvitation(ctx, i); err != nil {
		return model.ProjectInvitation{}, err
	}

	s.emailSender.SendEmailAsync(ctx, model.NewProjectInvitationEmail(p, tkn, i.Email))

	return i, nil
}

func (s *ProjectService) ListInvitations(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, err
	}

	invitations, err := s.repo.ListInvitationsForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if len(invitations) == 0 {
		exists, err := s.projectExists(ctx, projectID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apierrors.NewProjectNotFoundError()
		}
	}

	return invitations, nil
}

func (s *ProjectService) CancelInvitation(ctx context.Context, projectID, invitationID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return err
	}

	err := s.repo.DeleteInvitation(ctx, projectID, invitationID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) AcceptInvitation(ctx context.Context, tkn cryptox.Token) error {
	invitation, err := s.getPendingInvitationByToken(ctx, tkn)
	if err != nil {
		return err
	}

	u, err := s.userGetter.GetByEmail(ctx, invitation.Email)
	if err != nil && !apierrors.IsNotFoundError(err) {
		return err
	}

	// User does not exist yet, just accept the invitation
	// When a user registers, a project membership will be created;
	// PostgreSQL function handle_new_user() is triggered upon user creation
	if apierrors.IsNotFoundError(err) {
		invitation.Accept()
		return s.repo.UpdateInvitation(ctx, invitation)
	}

	// User exists
	member, err := model.NewProjectMember(u, invitation.ProjectID, invitation.ProjectRole)
	if err != nil {
		return err
	}

	return s.repo.CreateMember(ctx, member)
}

func (s *ProjectService) RejectInvitation(ctx context.Context, tkn cryptox.Token) error {
	// Only pending invitations can be rejected
	err := s.repo.DeleteInvitationByTokenHashAndStatus(ctx, tkn.ToHash(), model.InvitationStatusPending)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) ListMembers(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.ProjectMember, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, err
	}

	m, err := s.repo.ListMembersForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if len(m) == 0 {
		exists, err := s.projectExists(ctx, projectID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apierrors.NewProjectNotFoundError()
		}
	}

	return m, nil
}

func (s *ProjectService) DeleteMember(ctx context.Context, projectID, userID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return err
	}

	err := s.repo.DeleteMember(ctx, projectID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) UpdateMemberRole(ctx context.Context, newRole model.ProjectRole, projectID, userID, authUserID uuid.UUID) (model.ProjectMember, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.ProjectMember{}, err
	}

	exists, err := s.projectExists(ctx, projectID)
	if err != nil {
		return model.ProjectMember{}, err
	}
	if !exists {
		return model.ProjectMember{}, apierrors.NewProjectNotFoundError()
	}

	m, err := s.repo.ReadMember(ctx, projectID, userID)
	if err != nil {
		return model.ProjectMember{}, err
	}

	if err := m.UpdateProjectRole(newRole); err != nil {
		return model.ProjectMember{}, apierrors.NewProjectMemberUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.repo.UpdateMember(ctx, m); err != nil {
		return model.ProjectMember{}, err
	}

	return m, nil
}

func (s *ProjectService) projectExists(ctx context.Context, projectID uuid.UUID) (bool, error) {
	_, err := s.repo.ReadProject(ctx, projectID)
	if err != nil {
		if apierrors.IsNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *ProjectService) getPendingInvitationByToken(ctx context.Context, tkn cryptox.Token) (model.ProjectInvitation, error) {
	i, err := s.repo.ReadInvitationByTokenHashAndStatus(ctx, tkn.ToHash(), model.InvitationStatusPending)
	if err != nil {
		if dberrors.IsNotFoundError(err) {
			return model.ProjectInvitation{}, apierrors.NewProjectInvitationNotFoundError().Wrap(err)
		}

		return model.ProjectInvitation{}, err
	}

	return i, nil
}

func (s *ProjectService) invitationExists(ctx context.Context, email string, projectID uuid.UUID) (bool, error) {
	if _, err := s.repo.ReadInvitationByEmail(ctx, email, projectID); err != nil {
		if apierrors.IsNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *ProjectService) isEnvironmentNameUnique(ctx context.Context, projectID uuid.UUID, name string) (bool, error) {
	if _, err := s.repo.ReadEnvironmentByName(ctx, projectID, name); err != nil {
		if apierrors.IsNotFoundError(err) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

func (s *ProjectService) memberExists(ctx context.Context, projectID uuid.UUID, email string) (bool, error) {
	_, err := s.repo.ReadMemberByEmail(ctx, projectID, email)
	if err != nil {
		if apierrors.IsNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
