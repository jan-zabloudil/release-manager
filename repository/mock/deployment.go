package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type DeploymentRepository struct {
	mock.Mock
}

func (r *DeploymentRepository) Create(ctx context.Context, dpl svcmodel.Deployment) error {
	args := r.Called(ctx, dpl)
	return args.Error(0)
}
