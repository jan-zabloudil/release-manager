package service

import (
	"context"
	"fmt"

	svcerrors "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type DeploymentService struct {
	authGuard         authGuard
	projectGetter     projectGetter
	releaseGetter     releaseGetter
	environmentGetter environmentGetter
	repo              deploymentRepository
}

func NewDeploymentService(
	authGuard authGuard,
	projectGetter projectGetter,
	releaseGetter releaseGetter,
	environmentGetter environmentGetter,
	repo deploymentRepository,
) *DeploymentService {
	return &DeploymentService{
		authGuard:         authGuard,
		projectGetter:     projectGetter,
		releaseGetter:     releaseGetter,
		environmentGetter: environmentGetter,
		repo:              repo,
	}
}

func (s *DeploymentService) Create(
	ctx context.Context,
	input model.CreateDeploymentInput,
	projectID,
	authUserID uuid.UUID,
) (model.Deployment, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.Deployment{}, fmt.Errorf("authorizing project member: %w", err)
	}

	if err := input.Validate(); err != nil {
		return model.Deployment{}, svcerrors.NewDeploymentUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	exists, err := s.projectGetter.ProjectExists(ctx, projectID, authUserID)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("checking if project exists: %w", err)
	}
	if !exists {
		return model.Deployment{}, svcerrors.NewProjectNotFoundError()
	}

	// When getting release and env, project id is passed as a parameter to ensure that the objects belong to the project
	rls, err := s.releaseGetter.Get(ctx, projectID, input.ReleaseID, authUserID)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("getting release: %w", err)
	}
	env, err := s.environmentGetter.GetEnvironment(ctx, projectID, input.EnvironmentID, authUserID)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("getting environment: %w", err)
	}

	dpl := model.NewDeployment(rls, env, authUserID)

	if err := s.repo.Create(ctx, dpl); err != nil {
		return model.Deployment{}, fmt.Errorf("creating deployment: %w", err)
	}

	return dpl, nil
}

func (s *DeploymentService) ListForProject(
	ctx context.Context,
	projectID,
	authUserID uuid.UUID,
) ([]model.Deployment, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	// TODO add filtering options
	dpls, err := s.repo.ListForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing deployments: %w", err)
	}

	if len(dpls) == 0 {
		exists, err := s.projectGetter.ProjectExists(ctx, projectID, authUserID)
		if err != nil {
			return nil, fmt.Errorf("checking if project exists: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewProjectNotFoundError()
		}
	}

	return dpls, nil
}
