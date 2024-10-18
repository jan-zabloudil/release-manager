package mock

import (
	"context"

	"release-manager/pkg/id"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProjectService struct {
	mock.Mock
}

func (m *ProjectService) GetProject(ctx context.Context, projectID uuid.UUID, authUserID id.AuthUser) (model.Project, error) {
	args := m.Called(ctx, projectID, authUserID)
	return args.Get(0).(model.Project), args.Error(1)
}

func (m *ProjectService) GetEnvironment(ctx context.Context, projectID, envID uuid.UUID, authUserID id.AuthUser) (model.Environment, error) {
	args := m.Called(ctx, projectID, envID, authUserID)
	return args.Get(0).(model.Environment), args.Error(1)
}
