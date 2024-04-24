package mock

import (
	"context"

	cryptox "release-manager/pkg/crypto"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProjectInvitationRepository struct {
	mock.Mock
}

func (m *ProjectInvitationRepository) Create(ctx context.Context, i svcmodel.ProjectInvitation) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *ProjectInvitationRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectInvitationRepository) ReadByEmailForProject(ctx context.Context, email string, projectID uuid.UUID) (svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, email, projectID)
	return args.Get(0).(svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectInvitationRepository) ReadByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status svcmodel.ProjectInvitationStatus) (svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, hash, status)
	return args.Get(0).(svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectInvitationRepository) ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectInvitationRepository) Update(ctx context.Context, i svcmodel.ProjectInvitation) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *ProjectInvitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
