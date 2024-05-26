package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ReleaseService struct {
	mock.Mock
}

func (s *ReleaseService) Get(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) (svcmodel.Release, error) {
	args := s.Called(ctx, projectID, releaseID, authUserID)
	return args.Get(0).(svcmodel.Release), args.Error(1)
}
