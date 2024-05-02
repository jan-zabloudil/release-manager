package repository

import (
	"context"

	"release-manager/pkg/crypto"
	"release-manager/pkg/dberrors"
	"release-manager/repository/model"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type ProjectRepository struct {
	client            *supabase.Client
	projectEntity     string
	environmentEntity string
	invitationEntity  string
}

func NewProjectRepository(c *supabase.Client) *ProjectRepository {
	return &ProjectRepository{
		client:            c,
		projectEntity:     "projects",
		environmentEntity: "environments",
		invitationEntity:  "project_invitations",
	}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, p svcmodel.Project) error {
	data := model.ToProject(p)

	err := r.client.
		DB.From(r.projectEntity).
		Insert(&data).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) ReadProject(ctx context.Context, id uuid.UUID) (svcmodel.Project, error) {
	var resp model.Project
	err := r.client.
		DB.From(r.projectEntity).
		Select("*").Single().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Project{}, util.ToDBError(err)
	}

	p, err := model.ToSvcProject(resp)
	if err != nil {
		return svcmodel.Project{}, dberrors.NewToSvcModelError().Wrap(err)
	}

	return p, nil
}

func (r *ProjectRepository) ReadAllProjects(ctx context.Context) ([]svcmodel.Project, error) {
	var resp []model.Project
	err := r.client.
		DB.From(r.projectEntity).
		Select("*").
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	p, err := model.ToSvcProjects(resp)
	if err != nil {
		return nil, dberrors.NewToSvcModelError().Wrap(err)
	}

	return p, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	err := r.client.
		DB.From(r.projectEntity).
		Delete().Eq("id", id.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, p svcmodel.Project) error {
	data := model.ToUpdateProjectInput(p)

	err := r.client.
		DB.From(r.projectEntity).
		Update(&data).
		Eq("id", (p.ID).String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) CreateEnvironment(ctx context.Context, e svcmodel.Environment) error {
	data := model.ToEnvironment(e)

	err := r.client.
		DB.From(r.environmentEntity).
		Insert(&data).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) ReadEnvironment(ctx context.Context, envID uuid.UUID) (svcmodel.Environment, error) {
	var resp model.Environment
	err := r.client.
		DB.From(r.environmentEntity).
		Select("*").Single().
		Eq("id", envID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Environment{}, util.ToDBError(err)
	}

	env, err := model.ToSvcEnvironment(resp)
	if err != nil {
		return svcmodel.Environment{}, dberrors.NewToSvcModelError().Wrap(err)
	}

	return env, nil
}

func (r *ProjectRepository) ReadEnvironmentByNameForProject(ctx context.Context, projectID uuid.UUID, name string) (svcmodel.Environment, error) {
	var resp model.Environment
	err := r.client.
		DB.From(r.environmentEntity).
		Select("*").Single().
		Eq("name", name).
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.Environment{}, util.ToDBError(err)
	}

	env, err := model.ToSvcEnvironment(resp)
	if err != nil {
		return svcmodel.Environment{}, dberrors.NewToSvcModelError().Wrap(err)
	}

	return env, nil
}

func (r *ProjectRepository) ReadAllEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Environment, error) {
	var resp []model.Environment
	err := r.client.
		DB.From(r.environmentEntity).
		Select("*").
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	envs, err := model.ToSvcEnvironments(resp)
	if err != nil {
		return nil, dberrors.NewToSvcModelError().Wrap(err)
	}

	return envs, nil
}

func (r *ProjectRepository) DeleteEnvironment(ctx context.Context, envID uuid.UUID) error {
	err := r.client.
		DB.From(r.environmentEntity).
		Delete().
		Eq("id", envID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) UpdateEnvironment(ctx context.Context, e svcmodel.Environment) error {
	data := model.ToUpdateEnvironmentInput(e)

	err := r.client.
		DB.From(r.environmentEntity).
		Update(&data).
		Eq("id", e.ID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) CreateInvitation(ctx context.Context, i svcmodel.ProjectInvitation) error {
	data := model.ToProjectInvitation(i)

	err := r.client.
		DB.From(r.invitationEntity).
		Insert(&data).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) ReadInvitation(ctx context.Context, id uuid.UUID) (svcmodel.ProjectInvitation, error) {
	var resp model.ProjectInvitation
	err := r.client.
		DB.From(r.invitationEntity).
		Select("*").Single().
		Eq("id", id.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectInvitation{}, util.ToDBError(err)
	}

	return model.ToSvcProjectInvitation(resp), nil
}

func (r *ProjectRepository) ReadInvitationByEmailForProject(ctx context.Context, email string, projectID uuid.UUID) (svcmodel.ProjectInvitation, error) {
	var resp model.ProjectInvitation
	err := r.client.
		DB.From(r.invitationEntity).
		Select("*").Single().
		Eq("email", email).
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectInvitation{}, util.ToDBError(err)
	}

	return model.ToSvcProjectInvitation(resp), nil
}

func (r *ProjectRepository) ReadInvitationByTokenHashAndStatus(ctx context.Context, hash crypto.Hash, status svcmodel.ProjectInvitationStatus) (svcmodel.ProjectInvitation, error) {
	var resp model.ProjectInvitation
	err := r.client.
		DB.From(r.invitationEntity).
		Select("*").Single().
		Eq("token_hash", hash.ToBase64()).
		Eq("status", string(status)).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectInvitation{}, util.ToDBError(err)
	}

	return model.ToSvcProjectInvitation(resp), nil
}

func (r *ProjectRepository) ReadAllInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error) {
	var resp []model.ProjectInvitation
	err := r.client.
		DB.From(r.invitationEntity).
		Select("*").
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	return model.ToSvcProjectInvitations(resp), nil
}

func (r *ProjectRepository) UpdateInvitation(ctx context.Context, i svcmodel.ProjectInvitation) error {
	data := model.ToUpdateProjectInvitationInput(i)

	err := r.client.
		DB.From(r.invitationEntity).
		Update(&data).
		Eq("id", i.ID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) DeleteInvitation(ctx context.Context, id uuid.UUID) error {
	err := r.client.
		DB.From(r.invitationEntity).
		Delete().Eq("id", id.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}
