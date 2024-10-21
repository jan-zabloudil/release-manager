package mock

import (
	"context"

	"release-manager/pkg/crypto"
	"release-manager/pkg/id"
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

func (m *ProjectRepository) ListProjects(ctx context.Context) ([]svcmodel.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]svcmodel.Project), args.Error(1)
}

func (m *ProjectRepository) ListProjectsForUser(ctx context.Context, userID id.AuthUser) ([]svcmodel.Project, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]svcmodel.Project), args.Error(1)
}

func (m *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProjectRepository) UpdateProject(
	ctx context.Context,
	projectID uuid.UUID,
	updateFn func(p svcmodel.Project) (svcmodel.Project, error),
) error {
	args := m.Called(ctx, projectID, updateFn)
	return args.Error(0)
}

func (m *ProjectRepository) CreateEnvironment(ctx context.Context, e svcmodel.Environment) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *ProjectRepository) ReadEnvironment(ctx context.Context, projectID, envID uuid.UUID) (svcmodel.Environment, error) {
	args := m.Called(ctx, projectID, envID)
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

func (m *ProjectRepository) UpdateEnvironment(
	ctx context.Context,
	projectID,
	envID uuid.UUID,
	updateFn func(e svcmodel.Environment) (svcmodel.Environment, error),
) error {
	args := m.Called(ctx, projectID, envID, updateFn)
	return args.Error(0)
}

func (m *ProjectRepository) CreateInvitation(ctx context.Context, i svcmodel.ProjectInvitation) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *ProjectRepository) UpdateInvitation(
	ctx context.Context,
	invitationHash crypto.Hash,
	updateFn func(i svcmodel.ProjectInvitation) (svcmodel.ProjectInvitation, error),
) error {
	args := m.Called(ctx, invitationHash, updateFn)
	return args.Error(0)
}

func (m *ProjectRepository) ListInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.ProjectInvitation), args.Error(1)
}

func (m *ProjectRepository) DeleteInvitation(ctx context.Context, projectID, invitationID uuid.UUID) error {
	args := m.Called(ctx, projectID, invitationID)
	return args.Error(0)
}

func (m *ProjectRepository) DeleteInvitationByTokenHashAndStatus(ctx context.Context, hash crypto.Hash, status svcmodel.ProjectInvitationStatus) error {
	args := m.Called(ctx, hash, status)
	return args.Error(0)
}

func (m *ProjectRepository) CreateMember(
	ctx context.Context,
	tokenHash crypto.Hash,
	createMemberFn func(i svcmodel.ProjectInvitation) (svcmodel.ProjectMember, error),
) error {
	args := m.Called(ctx, tokenHash, createMemberFn)
	return args.Error(0)
}

func (m *ProjectRepository) ListMembersForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.ProjectMember), args.Error(1)
}

func (m *ProjectRepository) ListMembersForUser(ctx context.Context, userID id.AuthUser) ([]svcmodel.ProjectMember, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]svcmodel.ProjectMember), args.Error(1)
}

func (m *ProjectRepository) ReadMember(ctx context.Context, projectID uuid.UUID, userID id.User) (svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID)
	return args.Get(0).(svcmodel.ProjectMember), args.Error(1)
}

func (m *ProjectRepository) ReadMemberByEmail(ctx context.Context, projectID uuid.UUID, email string) (svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID, email)
	return args.Get(0).(svcmodel.ProjectMember), args.Error(1)
}

func (m *ProjectRepository) DeleteMember(ctx context.Context, projectID uuid.UUID, userID id.User) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *ProjectRepository) UpdateMember(
	ctx context.Context,
	projectID uuid.UUID,
	userID id.User,
	updateFn func(m svcmodel.ProjectMember) (svcmodel.ProjectMember, error),
) error {
	args := m.Called(ctx, projectID, userID, updateFn)
	return args.Error(0)
}
