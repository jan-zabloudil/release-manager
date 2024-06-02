package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type DeploymentRepository struct {
	mock.Mock
}

func (r *DeploymentRepository) Create(ctx context.Context, dpl svcmodel.Deployment) error {
	args := r.Called(ctx, dpl)
	return args.Error(0)
}

func (r *DeploymentRepository) ListForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Deployment, error) {
	args := r.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.Deployment), args.Error(1)
}
