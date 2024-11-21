package service

import (
	"context"
	"errors"
	"fmt"

	"release-manager/pkg/id"
	svcerrors "release-manager/service/errors"
	"release-manager/service/model"
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

func (s *ProjectService) CreateProject(ctx context.Context, input model.CreateProjectInput, authUserID id.AuthUser) (model.Project, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Project{}, fmt.Errorf("authorizing user role: %w", err)
	}

	// If empty release notification config is provided, use default config
	if input.ReleaseNotificationConfig.IsEmpty() {
		defaultCfg, err := s.getDefaultReleaseNotificationConfig(ctx)
		if err != nil {
			return model.Project{}, fmt.Errorf("getting default release notification config: %w", err)
		}

		input.ReleaseNotificationConfig = defaultCfg
	}

	p, err := model.NewProject(input)
	if err != nil {
		return model.Project{}, svcerrors.NewProjectInvalidError().Wrap(err).WithMessage(err.Error())
	}

	u, err := s.userGetter.GetAuthenticated(ctx, authUserID)
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

func (s *ProjectService) GetProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) (model.Project, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return model.Project{}, fmt.Errorf("authorizing project member: %w", err)
	}

	p, err := s.repo.ReadProject(ctx, projectID)
	if err != nil {
		return model.Project{}, fmt.Errorf("reading project: %w", err)
	}

	return p, err
}

func (s *ProjectService) ListProjects(ctx context.Context, authUserID id.AuthUser) ([]model.Project, error) {
	u, err := s.userGetter.GetAuthenticated(ctx, authUserID)
	if err != nil {
		return nil, fmt.Errorf("getting authenticated user: %w", err)
	}

	// Admin user can see all projects
	if u.IsAdmin() {
		p, err := s.repo.ListProjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing projects for admin user: %w", err)
		}

		return p, nil
	}

	p, err := s.repo.ListProjectsForUser(ctx, authUserID)
	if err != nil {
		return nil, fmt.Errorf("listing projects for non-admin user: %w", err)
	}

	return p, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return fmt.Errorf("authorizing user role: %w", err)
	}

	err := s.repo.DeleteProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}

	return nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, input model.UpdateProjectInput, projectID id.Project, authUserID id.AuthUser) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	if err := s.repo.UpdateProject(ctx, projectID, func(p model.Project) (model.Project, error) {
		if err := p.Update(input); err != nil {
			return model.Project{}, svcerrors.NewProjectInvalidError().Wrap(err).WithMessage(err.Error())
		}

		return p, nil
	}); err != nil {
		return fmt.Errorf("updating the project: %w", err)
	}

	return nil
}

func (s *ProjectService) SetGithubRepoForProject(ctx context.Context, rawRepoURL string, projectID id.Project, authUserID id.AuthUser) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return fmt.Errorf("getting Github token: %w", err)
	}

	repo, err := s.githubManager.ReadRepo(ctx, tkn, rawRepoURL)
	if err != nil {
		return fmt.Errorf("reading github repo: %w", err)
	}

	if err = s.repo.UpdateProject(ctx, projectID, func(p model.Project) (model.Project, error) {
		p.SetGithubRepo(&repo)
		return p, nil
	}); err != nil {
		return fmt.Errorf("updating project with Github repo: %w", err)
	}

	return nil
}

func (s *ProjectService) GetGithubRepoForProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) (model.GithubRepo, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.GithubRepo{}, fmt.Errorf("authorizing project member: %w", err)
	}

	p, err := s.repo.ReadProject(ctx, projectID)
	if err != nil {
		return model.GithubRepo{}, fmt.Errorf("reading project: %w", err)
	}

	if !p.IsGithubRepoSet() {
		return model.GithubRepo{}, svcerrors.NewGithubRepoNotSetForProjectError()
	}

	return *p.GithubRepo, nil
}

func (s *ProjectService) CreateEnvironment(ctx context.Context, input model.CreateEnvironmentInput, authUserID id.AuthUser) (model.Environment, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Environment{}, fmt.Errorf("authorizing user role: %w", err)
	}

	// Admin user was authorized (not project member), so we need to check if project exists
	exists, err := s.projectExists(ctx, input.ProjectID)
	if err != nil {
		return model.Environment{}, fmt.Errorf("checking if project exists: %w", err)
	}
	if !exists {
		return model.Environment{}, svcerrors.NewProjectNotFoundError()
	}

	env, err := model.NewEnvironment(input)
	if err != nil {
		return model.Environment{}, svcerrors.NewEnvironmentInvalidError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.repo.CreateEnvironment(ctx, env); err != nil {
		return model.Environment{}, fmt.Errorf("creating environment: %w", err)
	}

	return env, nil
}

func (s *ProjectService) GetEnvironment(ctx context.Context, projectID id.Project, envID id.Environment, authUserID id.AuthUser) (model.Environment, error) {
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
	input model.UpdateEnvironmentInput,
	projectID id.Project,
	envID id.Environment,
	authUserID id.AuthUser,
) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	if err := s.repo.UpdateEnvironment(ctx, projectID, envID, func(e model.Environment) (model.Environment, error) {
		if err := e.Update(input); err != nil {
			return model.Environment{}, svcerrors.NewEnvironmentInvalidError().Wrap(err).WithMessage(err.Error())
		}

		return e, nil
	}); err != nil {
		return fmt.Errorf("updating the environment: %w", err)
	}

	return nil
}

func (s *ProjectService) ListEnvironments(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]model.Environment, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	envs, err := s.repo.ListEnvironmentsForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing environments: %w", err)
	}

	return envs, nil
}

func (s *ProjectService) DeleteEnvironment(ctx context.Context, projectID id.Project, envID id.Environment, authUserID id.AuthUser) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return fmt.Errorf("authorizing user role: %w", err)
	}

	err := s.repo.DeleteEnvironment(ctx, projectID, envID)
	if err != nil {
		return fmt.Errorf("deleting environment: %w", err)
	}

	return nil
}

func (s *ProjectService) ListGithubRepoTags(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]model.GitTag, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Github token: %w", err)
	}

	p, err := s.repo.ReadProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("reading project: %w", err)
	}

	if !p.IsGithubRepoSet() {
		return nil, svcerrors.NewGithubRepoNotSetForProjectError()
	}

	t, err := s.githubManager.ReadTagsForRepo(ctx, tkn, *p.GithubRepo)
	if err != nil {
		return nil, fmt.Errorf("reading tags for github repo: %w", err)
	}

	return t, nil
}

func (s *ProjectService) Invite(ctx context.Context, input model.CreateProjectInvitationInput, authUserID id.AuthUser) (model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.ProjectInvitation{}, err
	}

	p, err := s.repo.ReadProject(ctx, input.ProjectID)
	if err != nil {
		return model.ProjectInvitation{}, fmt.Errorf("reading project: %w", err)
	}

	tkn, err := model.NewProjectInvitationToken()
	if err != nil {
		return model.ProjectInvitation{}, fmt.Errorf("creating token: %w", err)
	}

	i, err := model.NewProjectInvitation(input, tkn, authUserID)
	if err != nil {
		return model.ProjectInvitation{}, svcerrors.NewProjectInvitationInvalidError().Wrap(err).WithMessage(err.Error())
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

func (s *ProjectService) ListInvitations(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]model.ProjectInvitation, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}

	invitations, err := s.repo.ListInvitationsForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing invitations: %w", err)
	}

	// Admin user was authorized (not project member), so we need to check if project exists
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

func (s *ProjectService) CancelInvitation(ctx context.Context, projectID id.Project, invitationID id.ProjectInvitation, authUserID id.AuthUser) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return err
	}

	err := s.repo.DeleteInvitation(ctx, projectID, invitationID)
	if err != nil {
		return fmt.Errorf("deleting invitation: %w", err)
	}

	return nil
}

func (s *ProjectService) AcceptInvitation(ctx context.Context, tkn model.ProjectInvitationToken) error {
	// When accepting an invitation, two possible scenarios can happen:
	// 1. User does not exist yet, in which case invitation is only accepted.
	//    When a user registers later, a project member will be created (PostgreSQL function handle_new_user() is triggered upon user creation).
	//    Must be done via Postgres function because user registration is not controlled by this API.
	// 2. User already exists, in which case a project member is created and invitation is deleted.
	//
	// The current implementation aims to keep business logic in the service layer and prevent any race conditions.
	//
	// If we create a project member before updating or deleting the invitation, we might end up in an inconsistent state.
	// If the user does not exist when creating the member, the invitation would only be accepted.
	// However, if the user registers between these two steps (checking if the user exists and updating the invitation),
	// we could end up with an active user who has an accepted invitation but no associated project member.
	//
	// It is also possible to handle all logic within a single repository function, but that would require moving business logic into the repository layer.

	err := s.repo.UpdateInvitation(ctx, tkn.ToHash(), func(i model.ProjectInvitation) (model.ProjectInvitation, error) {
		if err := i.Accept(); err != nil {
			if errors.Is(err, model.ErrProjectInvitationAlreadyAccepted) {
				return model.ProjectInvitation{}, svcerrors.NewProjectInvitationNotFoundError().Wrap(err)
			}

			return model.ProjectInvitation{}, err
		}
		return i, nil
	})
	if err != nil {
		return fmt.Errorf("updating invitation: %w", err)
	}

	err = s.repo.CreateMember(ctx, tkn.ToHash(), func(i model.ProjectInvitation) (model.ProjectMember, error) {
		u, err := s.userGetter.GetByEmail(ctx, i.Email)
		if err != nil {
			return model.ProjectMember{}, fmt.Errorf("getting user by email: %w", err)
		}

		m, err := model.NewProjectMember(u, i.ProjectID, i.ProjectRole)
		if err != nil {
			return model.ProjectMember{}, fmt.Errorf("creating member object: %w", err)
		}

		return m, nil
	})
	if err != nil {
		if svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeUserNotFound) {
			// member is not created because user does not exist yet, that is expected case, not an error
			return nil
		}

		return fmt.Errorf("creating member: %w", err)
	}

	return nil
}

func (s *ProjectService) RejectInvitation(ctx context.Context, tkn model.ProjectInvitationToken) error {
	// Only pending invitations can be rejected
	if err := s.repo.DeleteInvitationByTokenHashAndStatus(ctx, tkn.ToHash(), model.InvitationStatusPending); err != nil {
		return fmt.Errorf("deleting invitation: %w", err)
	}

	return nil
}

func (s *ProjectService) ListMembersForProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]model.ProjectMember, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}

	m, err := s.repo.ListMembersForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing members for project: %w", err)
	}

	// Admin user was authorized (not project member), so we need to check if project exists
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

func (s *ProjectService) ListMembersForUser(ctx context.Context, authUserID id.AuthUser) ([]model.ProjectMember, error) {
	if err := s.authGuard.AuthorizeUserRoleUser(ctx, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}

	m, err := s.repo.ListMembersForUser(ctx, authUserID)
	if err != nil {
		return nil, fmt.Errorf("listing members for user: %w", err)
	}

	return m, nil
}

func (s *ProjectService) DeleteMember(ctx context.Context, projectID id.Project, userID id.User, authUserID id.AuthUser) error {
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
	projectID id.Project,
	userID id.User,
	authUserID id.AuthUser,
) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return fmt.Errorf("authorizing user role: %w", err)
	}

	if err := s.repo.UpdateMember(ctx, projectID, userID, func(m model.ProjectMember) (model.ProjectMember, error) {
		if err := m.UpdateProjectRole(newRole); err != nil {
			return model.ProjectMember{}, svcerrors.NewProjectMemberInvalidError().Wrap(err).WithMessage(err.Error())
		}

		return m, nil
	}); err != nil {
		return fmt.Errorf("updating member role: %w", err)
	}

	return nil
}

func (s *ProjectService) getDefaultReleaseNotificationConfig(ctx context.Context) (model.ReleaseNotificationConfig, error) {
	msg, err := s.settingsGetter.GetDefaultReleaseMessage(ctx)
	if err != nil {
		return model.ReleaseNotificationConfig{}, err
	}

	return model.ReleaseNotificationConfig{
		Message:            msg,
		ShowProjectName:    true,
		ShowReleaseTitle:   true,
		ShowReleaseNotes:   true,
		ShowLastDeployment: true,
		ShowSourceCode:     true,
	}, nil
}

func (s *ProjectService) projectExists(ctx context.Context, projectID id.Project) (bool, error) {
	if _, err := s.repo.ReadProject(ctx, projectID); err != nil {
		switch {
		case svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectNotFound):
			return false, nil
		default:
			return false, fmt.Errorf("reading project: %w", err)
		}
	}

	return true, nil
}

func (s *ProjectService) memberExists(ctx context.Context, projectID id.Project, email string) (bool, error) {
	if _, err := s.repo.ReadMemberByEmail(ctx, projectID, email); err != nil {
		switch {
		case svcerrors.IsErrorWithCode(err, svcerrors.ErrCodeProjectMemberNotFound):
			return false, nil
		default:
			return false, fmt.Errorf("reading member: %w", err)
		}
	}

	return true, nil
}
