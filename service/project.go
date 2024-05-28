package service

import (
	"context"
	"fmt"

	cryptox "release-manager/pkg/crypto"
	svcerrors "release-manager/service/errors"
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
		return model.Project{}, fmt.Errorf("authorizing user role: %w", err)
	}

	// If empty release notification config is provided, use default config
	if c.ReleaseNotificationConfig.IsEmpty() {
		defaultCfg, err := s.getDefaultReleaseNotificationConfig(ctx)
		if err != nil {
			return model.Project{}, fmt.Errorf("getting default release notification config: %w", err)
		}

		c.ReleaseNotificationConfig = defaultCfg
	}

	// TODO pass it to model.NewProject
	// TODO update also service update project method
	_, err := s.githubManager.GetGithubRepoFromRawURL(c.GithubRepositoryRawURL)
	if err != nil {
		return model.Project{}, fmt.Errorf("getting GitHub repo from URL: %w", err)
	}

	p, err := model.NewProject(c)
	if err != nil {
		return model.Project{}, svcerrors.NewProjectUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	u, err := s.userGetter.Get(ctx, authUserID, authUserID)
	if err != nil {
		return model.Project{}, fmt.Errorf("reading user: %w", err)
	}

	owner, err := model.NewProjectOwner(u, p.ID)
	if err != nil {
		return model.Project{}, fmt.Errorf("creating project owner object: %w", err)
	}

	if err := s.repo.CreateProjectWithOwner(ctx, p, owner); err != nil {
		return model.Project{}, fmt.Errorf("creating project and member in repository: %w", err)
	}

	return p, nil
}

func (s *ProjectService) GetProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (model.Project, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return model.Project{}, fmt.Errorf("authorizing project member: %w", err)
	}

	return s.getProject(ctx, projectID)
}

func (s *ProjectService) ProjectExists(ctx context.Context, projectID, authUserID uuid.UUID) (bool, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return false, fmt.Errorf("authorizing project member: %w", err)
	}

	return s.projectExists(ctx, projectID)
}

func (s *ProjectService) ListProjects(ctx context.Context, authUserID uuid.UUID) ([]model.Project, error) {
	err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID)
	switch {
	case err == nil:
		// Admin user can see all projects
		p, err := s.repo.ListProjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing projects for admin user: %w", err)
		}

		return p, nil
	case err != nil && svcerrors.IsInsufficientUserRoleError(err):
		// Non-admin user can see only projects they are members of
		p, err := s.repo.ListProjectsForUser(ctx, authUserID)
		if err != nil {
			return nil, fmt.Errorf("listing projects for non-admin user: %w", err)
		}

		return p, nil
	default:
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return fmt.Errorf("authorizing user role: %w", err)
	}

	err := s.repo.DeleteProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}

	return nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, u model.UpdateProjectInput, projectID, authUserID uuid.UUID) (model.Project, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.Project{}, fmt.Errorf("authorizing project member: %w", err)
	}

	p, err := s.repo.UpdateProject(ctx, projectID, func(p model.Project) (model.Project, error) {
		if err := p.Update(u); err != nil {
			return model.Project{}, svcerrors.NewProjectUnprocessableError().Wrap(err).WithMessage(err.Error())
		}

		return p, nil
	})
	if err != nil {
		return model.Project{}, fmt.Errorf("updating the project: %w", err)
	}

	return p, nil
}

func (s *ProjectService) CreateEnvironment(ctx context.Context, c model.CreateEnvironmentInput, authUserID uuid.UUID) (model.Environment, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Environment{}, fmt.Errorf("authorizing user role: %w", err)
	}

	_, err := s.getProject(ctx, c.ProjectID)
	if err != nil {
		return model.Environment{}, fmt.Errorf("reading project: %w", err)
	}

	env, err := model.NewEnvironment(c)
	if err != nil {
		return model.Environment{}, svcerrors.NewEnvironmentUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.repo.CreateEnvironment(ctx, env); err != nil {
		return model.Environment{}, fmt.Errorf("creating environment: %w", err)
	}

	return env, nil
}

func (s *ProjectService) GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (model.Environment, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return model.Environment{}, fmt.Errorf("authorizing project member: %w", err)
	}

	env, err := s.repo.ReadEnvironment(ctx, projectID, envID)
	if err != nil {
		return model.Environment{}, fmt.Errorf("reading environment: %w", err)
	}

	return env, nil
}

func (s *ProjectService) UpdateEnvironment(
	ctx context.Context,
	u model.UpdateEnvironmentInput,
	projectID,
	envID,
	authUserID uuid.UUID,
) (model.Environment, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.Environment{}, fmt.Errorf("authorizing project member: %w", err)
	}

	env, err := s.repo.UpdateEnvironment(ctx, projectID, envID, func(e model.Environment) (model.Environment, error) {
		if err := e.Update(u); err != nil {
			return model.Environment{}, svcerrors.NewEnvironmentUnprocessableError().Wrap(err).WithMessage(err.Error())
		}

		return e, nil
	})
	if err != nil {
		return model.Environment{}, fmt.Errorf("updating the environment: %w", err)
	}

	return env, nil
}

func (s *ProjectService) ListEnvironments(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.Environment, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	envs, err := s.repo.ListEnvironmentsForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing environments: %w", err)
	}

	if len(envs) == 0 {
		exists, err := s.projectExists(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("checking if project exists: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewProjectNotFoundError()
		}
	}

	return envs, nil
}

func (s *ProjectService) DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return fmt.Errorf("authorizing user role: %w", err)
	}

	err := s.repo.DeleteEnvironment(ctx, projectID, envID)
	if err != nil {
		return fmt.Errorf("deleting environment: %w", err)
	}

	return nil
}

func (s *ProjectService) ListGithubRepositoryTags(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.GitTag, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	project, err := s.getProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("reading project: %w", err)
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Github token: %w", err)
	}

	if !project.IsGithubConfigured() {
		return nil, svcerrors.NewGithubRepositoryNotConfiguredForProjectError()
	}

	t, err := s.githubManager.ReadTagsForRepository(ctx, tkn, project.GithubRepositoryURL)
	if err != nil {
		return nil, fmt.Errorf("reading tags for github repository: %w", err)
	}

	return t, nil
}

func (s *ProjectService) Invite(ctx context.Context, c model.CreateProjectInvitationInput, authUserID uuid.UUID) (model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.ProjectInvitation{}, err
	}

	p, err := s.getProject(ctx, c.ProjectID)
	if err != nil {
		return model.ProjectInvitation{}, fmt.Errorf("reading project: %w", err)
	}

	tkn, err := cryptox.NewToken()
	if err != nil {
		return model.ProjectInvitation{}, fmt.Errorf("creating token: %w", err)
	}

	i, err := model.NewProjectInvitation(c, tkn, authUserID)
	if err != nil {
		return model.ProjectInvitation{}, svcerrors.NewProjectInvitationUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	memberExists, err := s.memberExists(ctx, i.ProjectID, i.Email)
	if err != nil {
		return model.ProjectInvitation{}, fmt.Errorf("checking if member exists: %w", err)
	}
	if memberExists {
		return model.ProjectInvitation{}, svcerrors.NewProjectMemberAlreadyExistsError()
	}

	if err := s.repo.CreateInvitation(ctx, i); err != nil {
		return model.ProjectInvitation{}, fmt.Errorf("creating invitation: %w", err)
	}

	s.emailSender.SendProjectInvitationEmailAsync(ctx, model.NewProjectInvitationEmailData(p.Name, tkn), i.Email)

	return i, nil
}

func (s *ProjectService) ListInvitations(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}

	invitations, err := s.repo.ListInvitationsForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing invitations: %w", err)
	}

	if len(invitations) == 0 {
		exists, err := s.projectExists(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("checking if project exists: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewProjectNotFoundError()
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
		return fmt.Errorf("deleting invitation: %w", err)
	}

	return nil
}

func (s *ProjectService) AcceptInvitation(ctx context.Context, tkn cryptox.Token) error {
	// When accepting an invitation, two possible scenarios can happen:
	// 1. User does not exist yet, in which case invitation is only accepted
	//    When a user registers later, a project member will be created (PostgreSQL function handle_new_user() is triggered upon user creation)
	//    Must be done via Postgres function because user registration is not controlled by this API
	// 2. User already exists, in which case a project member is created and invitation is deleted
	//
	// To keep more logic in service, we first fetch the invitation to get invitation's email
	// Then we check if user with given email exists
	// First two steps could be done in one query, but it would require mixing entities from two different repositories
	// Based on user existence we either accept the invitation or create a project member and delete the invitation
	//
	// There is one possible race condition:
	// 1. We check if user exists (and he does not)
	// 2. We only accept the invitation
	// But user would register between these two steps
	// Resulting in a project member not being created for existing user
	// The race condition is handled by the PostgreSQL function check_accepted_invitations_for_registered_users()

	invitation, err := s.repo.ReadPendingInvitationByHash(ctx, tkn.ToHash())
	if err != nil {
		return fmt.Errorf("reading invitation: %w", err)
	}

	u, err := s.userGetter.GetByEmail(ctx, invitation.Email)
	if err != nil && !svcerrors.IsNotFoundError(err) {
		return fmt.Errorf("reading user: %w", err)
	}

	// User does not exist yet
	if svcerrors.IsNotFoundError(err) {
		if err := s.repo.AcceptPendingInvitation(ctx, invitation.ID, func(i *model.ProjectInvitation) {
			i.Accept()
		}); err != nil {
			return fmt.Errorf("accepting invitation: %w", err)
		}

		return nil
	}

	member, err := model.NewProjectMember(u, invitation.ProjectID, invitation.ProjectRole)
	if err != nil {
		return fmt.Errorf("creating project member object: %w", err)
	}

	// Creates a project member and deletes the invitation
	if err := s.repo.CreateMember(ctx, member); err != nil {
		return fmt.Errorf("creating project member in repo: %w", err)
	}

	return nil
}

func (s *ProjectService) RejectInvitation(ctx context.Context, tkn cryptox.Token) error {
	// Only pending invitations can be rejected
	if err := s.repo.DeleteInvitationByTokenHashAndStatus(ctx, tkn.ToHash(), model.InvitationStatusPending); err != nil {
		return fmt.Errorf("deleting invitation: %w", err)
	}

	return nil
}

func (s *ProjectService) ListMembers(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.ProjectMember, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}

	m, err := s.repo.ListMembersForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing members: %w", err)
	}

	if len(m) == 0 {
		exists, err := s.projectExists(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("checking if project exists: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewProjectNotFoundError()
		}
	}

	return m, nil
}

func (s *ProjectService) DeleteMember(ctx context.Context, projectID, userID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return fmt.Errorf("authorizing user role: %w", err)
	}

	if err := s.repo.DeleteMember(ctx, projectID, userID); err != nil {
		return fmt.Errorf("deleting member: %w", err)
	}

	return nil
}

func (s *ProjectService) UpdateMemberRole(
	ctx context.Context,
	newRole model.ProjectRole,
	projectID,
	userID, authUserID uuid.UUID,
) (model.ProjectMember, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.ProjectMember{}, fmt.Errorf("authorizing user role: %w", err)
	}

	m, err := s.repo.UpdateMemberRole(ctx, projectID, userID, func(m model.ProjectMember) (model.ProjectMember, error) {
		if err := m.UpdateProjectRole(newRole); err != nil {
			return model.ProjectMember{}, svcerrors.NewProjectMemberUnprocessableError().Wrap(err).WithMessage(err.Error())
		}

		return m, nil
	})
	if err != nil {
		return model.ProjectMember{}, fmt.Errorf("updating member role: %w", err)
	}

	return m, nil
}

func (s *ProjectService) getProject(ctx context.Context, projectID uuid.UUID) (model.Project, error) {
	p, err := s.repo.ReadProject(ctx, projectID)
	if err != nil {
		return model.Project{}, fmt.Errorf("reading project: %w", err)
	}

	return p, nil
}

func (s *ProjectService) getDefaultReleaseNotificationConfig(ctx context.Context) (model.ReleaseNotificationConfig, error) {
	msg, err := s.settingsGetter.GetDefaultReleaseMessage(ctx)
	if err != nil {
		return model.ReleaseNotificationConfig{}, err
	}

	return model.ReleaseNotificationConfig{
		Message:          msg,
		ShowProjectName:  true,
		ShowReleaseTitle: true,
		ShowReleaseNotes: true,
		ShowDeployments:  true,
		ShowSourceCode:   true,
	}, nil
}

func (s *ProjectService) projectExists(ctx context.Context, projectID uuid.UUID) (bool, error) {
	_, err := s.repo.ReadProject(ctx, projectID)
	if err != nil {
		if svcerrors.IsNotFoundError(err) {
			return false, nil
		}

		return false, fmt.Errorf("reading project: %w", err)
	}

	return true, nil
}

func (s *ProjectService) memberExists(ctx context.Context, projectID uuid.UUID, email string) (bool, error) {
	_, err := s.repo.ReadMemberByEmail(ctx, projectID, email)
	if err != nil {
		if svcerrors.IsNotFoundError(err) {
			return false, nil
		}

		return false, fmt.Errorf("reading member: %w", err)
	}

	return true, nil
}
