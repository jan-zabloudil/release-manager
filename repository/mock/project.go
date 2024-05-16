package mock

import (
	"context"

	"release-manager/pkg/crypto"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProjectRepository struct {
	mock.Mock
}

func (m *ProjectRepository) CreateProjectWithOwner(ctx context.Context, p svcmodel.Project, owner svcmodel.ProjectMember) error {
	args := m.Called(ctx, p, owner)
	return args.Error(0)
}

func (m *ProjectRepository) ReadProject(ctx context.Context, id uuid.UUID) (svcmodel.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(svcmodel.Project), args.Error(1)
}

func (m *ProjectRepository) ReadAllProjects(ctx context.Context) ([]svcmodel.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]svcmodel.Project), args.Error(1)
}

func (m *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProjectRepository) UpdateProject(ctx context.Context, projectID uuid.UUID, fn svcmodel.UpdateProjectFunc) (svcmodel.Project, error) {
	args := m.Called(ctx, projectID, fn)
	return args.Get(0).(svcmodel.Project), args.Error(1)
}

func (m *ProjectRepository) CreateEnvironment(ctx context.Context, e svcmodel.Environment) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *ProjectRepository) ReadEnvironment(ctx context.Context, projectID, envID uuid.UUID) (svcmodel.Environment, error) {
	args := m.Called(ctx, projectID, envID)
	return args.Get(0).(svcmodel.Environment), args.Error(1)
}

func (m *ProjectRepository) ReadEnvironmentByName(ctx context.Context, projectID uuid.UUID, name string) (svcmodel.Environment, error) {
	args := m.Called(ctx, projectID, name)
	return args.Get(0).(svcmodel.Environment), args.Error(1)
}

func (m *ProjectRepository) ListEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Environment, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.Environment), args.Error(1)
}

func (m *ProjectRepository) DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID) error {
	args := m.Called(ctx, projectID, envID)
	return args.Error(0)
}

func (m *ProjectRepository) UpdateEnvironment(ctx context.Context, e svcmodel.Environment) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *ProjectRepository) CreateInvitation(ctx context.Context, i svcmodel.ProjectInvitation) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *ProjectRepository) ReadInvitationByEmailForProject(ctx context.Context, email string, projectID uuid.UUID) (svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, email, projectID)
	return args.Get(0).(svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectRepository) ReadInvitationByTokenHashAndStatus(ctx context.Context, hash crypto.Hash, status svcmodel.ProjectInvitationStatus) (svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, hash, status)
	return args.Get(0).(svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectRepository) ReadAllInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectRepository) UpdateInvitation(ctx context.Context, i svcmodel.ProjectInvitation) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *ProjectRepository) DeleteInvitation(ctx context.Context, projectID, invitationID uuid.UUID) error {
	args := m.Called(ctx, projectID, invitationID)
	return args.Error(0)
}

func (m *ProjectRepository) DeleteInvitationByTokenHashAndStatus(ctx context.Context, hash crypto.Hash, status svcmodel.ProjectInvitationStatus) error {
	args := m.Called(ctx, hash, status)
	return args.Error(0)
}

func (m *ProjectRepository) CreateMember(ctx context.Context, member svcmodel.ProjectMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *ProjectRepository) ListMembersForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.ProjectMember), args.Error(1)
}

func (m *ProjectRepository) ReadMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID)
	return args.Get(0).(svcmodel.ProjectMember), args.Error(1)
}

func (m *ProjectRepository) DeleteMember(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *ProjectRepository) UpdateMember(ctx context.Context, pm svcmodel.ProjectMember) error {
	args := m.Called(ctx, pm)
	return args.Error(0)
}
