package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProjectService struct {
	mock.Mock
}

func (m *ProjectService) Get(ctx context.Context, projectID, authUserID uuid.UUID) (model.Project, error) {
	args := m.Called(ctx, projectID, authUserID)
	return args.Get(0).(model.Project), args.Error(1)
}
