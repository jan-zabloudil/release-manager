package repository

import (
	"context"
	"errors"
	"fmt"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/crypto"
	"release-manager/pkg/dberrors"
	"release-manager/repository/model"
	"release-manager/repository/query"
	"release-manager/repository/util"
	svcmodel "release-manager/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"
)

const (
	projectDBEntity       = "projects"
	environmentDBEntity   = "environments"
	invitationDBEntity    = "project_invitations"
	projectMemberDBEntity = "project_members"
)

type ProjectRepository struct {
	client *supabase.Client
	dbpool *pgxpool.Pool
}

func NewProjectRepository(c *supabase.Client, pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		client: c,
		dbpool: pool,
	}
}

func (r *ProjectRepository) CreateProjectWithOwner(ctx context.Context, p svcmodel.Project, owner svcmodel.ProjectMember) (err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, query.CreateProject, pgx.NamedArgs{
		"id":                        p.ID,
		"name":                      p.Name,
		"slackChannelID":            p.SlackChannelID,
		"releaseNotificationConfig": model.ReleaseNotificationConfig(p.ReleaseNotificationConfig), // converted to the struct with json tags (the field is saved as json in the database)
		"createdAt":                 p.CreatedAt,
		"updatedAt":                 p.UpdatedAt,
		"githubRepositoryURL":       p.GithubRepositoryURL.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	_, err = tx.Exec(ctx, query.CreateProjectMember, pgx.NamedArgs{
		"userID":      owner.User.ID,
		"projectID":   p.ID,
		"projectRole": owner.ProjectRole,
		"createdAt":   owner.CreatedAt,
		"updatedAt":   owner.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create project member: %w", err)
	}

	return nil
}

func (r *ProjectRepository) ReadProject(ctx context.Context, id uuid.UUID) (svcmodel.Project, error) {
	var p model.Project

	err := pgxscan.Get(ctx, r.dbpool, &p, query.ReadProject, pgx.NamedArgs{"id": id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.Project{}, apierrors.NewProjectNotFoundError().Wrap(err)
		}

		return svcmodel.Project{}, err
	}

	return model.ToSvcProject(p)
}

func (r *ProjectRepository) ReadAllProjects(ctx context.Context) ([]svcmodel.Project, error) {
	var resp []model.Project
	err := r.client.
		DB.From(projectDBEntity).
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
	result, err := r.dbpool.Exec(ctx, query.DeleteProject, pgx.NamedArgs{"id": id})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apierrors.NewProjectNotFoundError().Wrap(err)
	}

	return nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, p svcmodel.Project) error {
	data := model.ToUpdateProjectInput(p)

	err := r.client.
		DB.From(projectDBEntity).
		Update(&data).
		Eq("id", (p.ID).String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

func (r *ProjectRepository) CreateEnvironment(ctx context.Context, e svcmodel.Environment) error {
	_, err := r.dbpool.Exec(ctx, query.CreateEnvironment, pgx.NamedArgs{
		"id":         e.ID,
		"projectID":  e.ProjectID,
		"name":       e.Name,
		"serviceURL": e.ServiceURL.String(),
		"createdAt":  e.CreatedAt,
		"updatedAt":  e.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) ReadEnvironment(ctx context.Context, envID uuid.UUID) (svcmodel.Environment, error) {
	var resp model.Environment
	err := r.client.
		DB.From(environmentDBEntity).
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
		DB.From(environmentDBEntity).
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
		DB.From(environmentDBEntity).
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
		DB.From(environmentDBEntity).
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
		DB.From(environmentDBEntity).
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
		DB.From(invitationDBEntity).
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
		DB.From(invitationDBEntity).
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
		DB.From(invitationDBEntity).
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
		DB.From(invitationDBEntity).
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
		DB.From(invitationDBEntity).
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
		DB.From(invitationDBEntity).
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
		DB.From(invitationDBEntity).
		Delete().Eq("id", id.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}

// CreateMember creates a project member and deletes the invitation
func (r *ProjectRepository) CreateMember(ctx context.Context, m svcmodel.ProjectMember) (err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, query.CreateProjectMember, pgx.NamedArgs{
		"userID":      m.User.ID,
		"projectID":   m.ProjectID,
		"projectRole": m.ProjectRole,
		"createdAt":   m.CreatedAt,
		"updatedAt":   m.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create project member: %w", err)
	}

	_, err = tx.Exec(ctx, query.DeleteProjectInvitationByEmailAndProjectID, pgx.NamedArgs{
		"email":     m.User.Email,
		"projectID": m.ProjectID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete project invitation: %w", err)
	}

	return nil
}

func (r *ProjectRepository) ReadMembersForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error) {
	var resp []model.ProjectMember
	err := r.client.
		DB.From(projectMemberDBEntity).
		Select(fmt.Sprintf("*,%s(*)", userDBEntity)). // docs https://supabase.com/docs/guides/database/joins-and-nesting
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, util.ToDBError(err)
	}

	return model.ToSvcProjectMembers(resp), nil
}

func (r *ProjectRepository) ReadMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (svcmodel.ProjectMember, error) {
	var resp model.ProjectMember
	err := r.client.
		DB.From(projectMemberDBEntity).
		Select(fmt.Sprintf("*,%s(*)", userDBEntity)). // docs https://supabase.com/docs/guides/database/joins-and-nesting
		Single().
		Eq("project_id", projectID.String()).
		Eq("user_id", userID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectMember{}, util.ToDBError(err)
	}

	return model.ToSvcProjectMember(resp), nil
}

func (r *ProjectRepository) DeleteMember(ctx context.Context, projectID, userID uuid.UUID) error {
	result, err := r.dbpool.Exec(ctx, query.DeleteProjectMember, pgx.NamedArgs{
		"projectID": projectID,
		"userID":    userID,
	})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apierrors.NewProjectMemberNotFoundError()
	}

	return nil
}

func (r *ProjectRepository) UpdateMember(ctx context.Context, m svcmodel.ProjectMember) error {
	data := model.ToUpdateProjectMemberInput(m)

	err := r.client.
		DB.From(projectMemberDBEntity).
		Update(&data).
		Eq("project_id", m.ProjectID.String()).
		Eq("user_id", m.User.ID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return util.ToDBError(err)
	}

	return nil
}
