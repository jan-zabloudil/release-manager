package service

import (
	"context"
	"fmt"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/dberrors"
	"release-manager/pkg/githuberrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectService struct {
	authSvc        model.AuthService
	settingsSvc    model.SettingsService
	projectRepo    model.ProjectRepository
	envRepo        model.EnvironmentRepository
	invitationRepo model.ProjectInvitationRepository
	githubClient   model.GithubClient
}

func NewProjectService(
	authSvc model.AuthService,
	settingsSvc model.SettingsService,
	projectRepo model.ProjectRepository,
	envRepo model.EnvironmentRepository,
	invitationRepo model.ProjectInvitationRepository,
	githubClient model.GithubClient,
) *ProjectService {
	return &ProjectService{
		authSvc:        authSvc,
		settingsSvc:    settingsSvc,
		projectRepo:    projectRepo,
		envRepo:        envRepo,
		invitationRepo: invitationRepo,
		githubClient:   githubClient,
	}
}

func (s *ProjectService) Create(ctx context.Context, c model.CreateProjectInput, authUserID uuid.UUID) (model.Project, error) {
	if err := s.authSvc.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return model.Project{}, err
	}

	p, err := model.NewProject(c)
	if err != nil {
		return model.Project{}, apierrors.NewProjectUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.projectRepo.Create(ctx, p); err != nil {
		return model.Project{}, err
	}

	return p, nil
}

func (s *ProjectService) Get(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (model.Project, error) {
	// TODO add project member authorization

	p, err := s.projectRepo.Read(ctx, projectID)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return model.Project{}, apierrors.NewProjectNotFoundError().Wrap(err)
		default:
			return model.Project{}, err
		}
	}

	return p, nil
}

func (s *ProjectService) ListAll(ctx context.Context, authUserID uuid.UUID) ([]model.Project, error) {
	// TODO add project member authorization
	// TODO fetch only project where the user is a member

	return s.projectRepo.ReadAll(ctx)
}

func (s *ProjectService) Delete(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authSvc.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return err
	}

	_, err := s.Get(ctx, projectID, authUserID)
	if err != nil {
		return err
	}

	return s.projectRepo.Delete(ctx, projectID)
}

func (s *ProjectService) Update(ctx context.Context, u model.UpdateProjectInput, projectID, authUserID uuid.UUID) (model.Project, error) {
	// TODO add project member authorization

	p, err := s.Get(ctx, projectID, authUserID)
	if err != nil {
		return model.Project{}, err
	}

	if err := p.Update(u); err != nil {
		return model.Project{}, apierrors.NewProjectUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.projectRepo.Update(ctx, p); err != nil {
		return model.Project{}, err
	}

	return p, nil
}

func (s *ProjectService) CreateEnvironment(ctx context.Context, c model.CreateEnvironmentInput, authUserID uuid.UUID) (model.Environment, error) {
	if err := s.authSvc.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return model.Environment{}, err
	}

	_, err := s.Get(ctx, c.ProjectID, authUserID)
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

	if err := s.envRepo.Create(ctx, env); err != nil {
		return model.Environment{}, err
	}

	return env, nil
}

func (s *ProjectService) GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (model.Environment, error) {
	// TODO add project member authorization

	_, err := s.Get(ctx, projectID, authUserID)
	if err != nil {
		return model.Environment{}, err
	}

	env, err := s.envRepo.Read(ctx, envID)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return model.Environment{}, apierrors.NewEnvironmentNotFoundError().Wrap(err)
		default:
			return model.Environment{}, err
		}
	}

	return env, nil
}

func (s *ProjectService) UpdateEnvironment(ctx context.Context, u model.UpdateEnvironmentInput, projectID, envID, authUserID uuid.UUID) (model.Environment, error) {
	// TODO add project member authorization

	if _, err := s.Get(ctx, projectID, authUserID); err != nil {
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

	if err := s.envRepo.Update(ctx, env); err != nil {
		return model.Environment{}, err
	}

	return env, nil
}

func (s *ProjectService) ListEnvironments(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.Environment, error) {
	// TODO add project member authorization

	_, err := s.Get(ctx, projectID, authUserID)
	if err != nil {
		return nil, err
	}

	envs, err := s.envRepo.ReadAllForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return envs, nil
}

func (s *ProjectService) DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authSvc.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return err
	}

	_, err := s.Get(ctx, projectID, authUserID)
	if err != nil {
		return err
	}

	_, err = s.GetEnvironment(ctx, projectID, envID, authUserID)
	if err != nil {
		return err
	}

	return s.envRepo.Delete(ctx, envID)
}

func (s *ProjectService) ListGithubRepositoryTags(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.GitTag, error) {
	// TODO add project member authorization

	project, err := s.Get(ctx, projectID, authUserID)
	if err != nil {
		return nil, err
	}

	github, err := s.settingsSvc.GetGithubSettings(ctx)
	if err != nil {
		return nil, err
	}
	if !github.Enabled {
		return nil, apierrors.NewGithubIntegrationNotEnabledError()
	}

	if !project.IsGithubConfigured() {
		return nil, apierrors.NewGithubRepositoryNotConfiguredForProjectError()
	}

	s.githubClient.SetToken(github.Token)

	t, err := s.githubClient.ListTagsForRepository(ctx, project.GithubRepository)
	if err != nil {
		switch {
		case githuberrors.IsUnauthorizedError(err):
			return nil, apierrors.NewGithubIntegrationUnauthorizedError().Wrap(err)
		case githuberrors.IsForbiddenError(err):
			return nil, apierrors.NewGithubIntegrationForbiddenError().Wrap(err)
		case githuberrors.IsNotFoundError(err):
			return nil, apierrors.NewGithubRepositoryNotFoundError().Wrap(err).
				WithMessage(fmt.Sprintf("%s not found among accessible repositories.", project.GithubRepository.URL.String()))
		default:
			return nil, err
		}
	}

	return t, nil
}

func (s *ProjectService) isEnvironmentNameUnique(ctx context.Context, projectID uuid.UUID, name string) (bool, error) {
	if _, err := s.envRepo.ReadByNameForProject(ctx, projectID, name); err != nil {
		if dberrors.IsNotFoundError(err) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
