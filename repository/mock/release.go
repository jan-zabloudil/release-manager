package mock

import (
	"context"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type ReleaseRepository struct {
	mock.Mock
}

func (m *ReleaseRepository) CreateRelease(ctx context.Context, rls svcmodel.Release) error {
	args := m.Called(ctx, rls)
	return args.Error(0)
}

func (m *ReleaseRepository) ReadRelease(ctx context.Context, releaseID id.Release) (svcmodel.Release, error) {
	args := m.Called(ctx, releaseID)
	return args.Get(0).(svcmodel.Release), args.Error(1)
}

func (m *ReleaseRepository) ReadReleaseForProject(ctx context.Context, projectID id.Project, releaseID id.Release) (svcmodel.Release, error) {
	args := m.Called(ctx, projectID, releaseID)
	return args.Get(0).(svcmodel.Release), args.Error(1)
}

func (m *ReleaseRepository) DeleteRelease(ctx context.Context, releaseID id.Release) error {
	args := m.Called(ctx, releaseID)
	return args.Error(0)
}

func (m *ReleaseRepository) ListReleasesForProject(ctx context.Context, projectID id.Project) ([]svcmodel.Release, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.Release), args.Error(1)
}

func (m *ReleaseRepository) UpdateRelease(
	ctx context.Context,
	releaseID id.Release,
	updateFn func(r svcmodel.Release) (svcmodel.Release, error),
) error {
	args := m.Called(ctx, releaseID, updateFn)
	return args.Error(0)
}

func (m *ReleaseRepository) CreateDeployment(ctx context.Context, dpl svcmodel.Deployment) error {
	args := m.Called(ctx, dpl)
	return args.Error(0)
}

func (m *ReleaseRepository) ListDeploymentsForProject(ctx context.Context, params svcmodel.ListDeploymentsFilterParams, projectID id.Project) ([]svcmodel.Deployment, error) {
	args := m.Called(ctx, params, projectID)
	return args.Get(0).([]svcmodel.Deployment), args.Error(1)
}

func (m *ReleaseRepository) ReadLastDeploymentForRelease(ctx context.Context, releaseID id.Release) (svcmodel.Deployment, error) {
	args := m.Called(ctx, releaseID)
	return args.Get(0).(svcmodel.Deployment), args.Error(1)
}

func (m *ReleaseRepository) DeleteReleaseByGitTag(ctx context.Context, repo svcmodel.GithubRepo, tagName string) error {
	args := m.Called(ctx, repo, tagName)
	return args.Error(0)
}
