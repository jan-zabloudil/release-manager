package repository

import (
	"context"
	"errors"
	"fmt"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/crypto"
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

	_, err = tx.Exec(ctx, query.CreateMember, pgx.NamedArgs{
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
	return r.readProject(ctx, r.dbpool, query.ReadProject, id)
}

func (r *ProjectRepository) readProject(ctx context.Context, q querier, query string, id uuid.UUID) (svcmodel.Project, error) {
	var p model.Project

	err := pgxscan.Get(ctx, q, &p, query, pgx.NamedArgs{"id": id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.Project{}, apierrors.NewProjectNotFoundError().Wrap(err)
		}

		return svcmodel.Project{}, err
	}

	return model.ToSvcProject(p)
}

func (r *ProjectRepository) ListProjects(ctx context.Context) ([]svcmodel.Project, error) {
	return r.listProjects(ctx, query.ListProjects, nil)
}

func (r *ProjectRepository) ListProjectsForUser(ctx context.Context, userID uuid.UUID) ([]svcmodel.Project, error) {
	return r.listProjects(ctx, query.ListProjectsForUser, pgx.NamedArgs{"userID": userID})
}

func (r *ProjectRepository) listProjects(ctx context.Context, readQuery string, args pgx.NamedArgs) ([]svcmodel.Project, error) {
	var p []model.Project

	err := pgxscan.Select(ctx, r.dbpool, &p, readQuery, args)
	if err != nil {
		return nil, err
	}

	return model.ToSvcProjects(p)
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	result, err := r.dbpool.Exec(ctx, query.DeleteProject, pgx.NamedArgs{"id": id})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apierrors.NewProjectNotFoundError()
	}

	return nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, projectID uuid.UUID, fn svcmodel.UpdateProjectFunc) (p svcmodel.Project, err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return svcmodel.Project{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	p, err = r.readProject(ctx, tx, query.AppendForUpdate(query.ReadProject), projectID)
	if err != nil {
		return svcmodel.Project{}, fmt.Errorf("failed to read project: %w", err)
	}

	p, err = fn(p)
	if err != nil {
		return svcmodel.Project{}, err
	}

	_, err = tx.Exec(ctx, query.UpdateProject, pgx.NamedArgs{
		"id":                        p.ID,
		"name":                      p.Name,
		"slackChannelID":            p.SlackChannelID,
		"releaseNotificationConfig": model.ReleaseNotificationConfig(p.ReleaseNotificationConfig), // converted to the struct with json tags (the field is saved as json in the database)
		"githubRepositoryURL":       p.GithubRepositoryURL.String(),
		"updatedAt":                 p.UpdatedAt,
	})
	if err != nil {
		return svcmodel.Project{}, fmt.Errorf("failed to update project: %w", err)
	}

	return p, nil
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

func (r *ProjectRepository) ReadEnvironment(ctx context.Context, projectID, envID uuid.UUID) (svcmodel.Environment, error) {
	// Project ID is not needed in the query because envID is primary key
	// But it is added for security reasons
	// To make sure that the environment belongs to the project that is passed from the service
	return r.readEnvironment(ctx, r.dbpool, query.ReadEnvironment, pgx.NamedArgs{
		"envID":     envID,
		"projectID": projectID,
	})
}

func (r *ProjectRepository) ReadEnvironmentByName(ctx context.Context, projectID uuid.UUID, name string) (svcmodel.Environment, error) {
	// Fetches the environment by name for the project
	return r.readEnvironment(ctx, r.dbpool, query.ReadEnvironmentByName, pgx.NamedArgs{
		"name":      name,
		"projectID": projectID,
	})
}

func (r *ProjectRepository) readEnvironment(ctx context.Context, q querier, readQuery string, args pgx.NamedArgs) (svcmodel.Environment, error) {
	var e model.Environment

	err := pgxscan.Get(ctx, q, &e, readQuery, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.Environment{}, apierrors.NewEnvironmentNotFoundError().Wrap(err)
		}

		return svcmodel.Environment{}, err
	}

	return model.ToSvcEnvironment(e)
}

func (r *ProjectRepository) ListEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Environment, error) {
	var e []model.Environment

	err := pgxscan.Select(ctx, r.dbpool, &e, query.ListEnvironmentsForProject, pgx.NamedArgs{"projectID": projectID})
	if err != nil {
		return nil, err
	}

	return model.ToSvcEnvironments(e)
}

func (r *ProjectRepository) DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID) error {
	result, err := r.dbpool.Exec(ctx, query.DeleteEnvironment, pgx.NamedArgs{
		"envID":     envID,
		"projectID": projectID,
	})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apierrors.NewEnvironmentNotFoundError()
	}

	return nil
}

func (r *ProjectRepository) UpdateEnvironment(
	ctx context.Context,
	projectID,
	envID uuid.UUID,
	fn svcmodel.UpdateEnvironmentFunc,
) (env svcmodel.Environment, err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return svcmodel.Environment{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	// Project ID is not needed in the query because envID is primary key
	// But it is added for security reasons
	// To make sure that the environment belongs to the project that is passed from the service
	env, err = r.readEnvironment(ctx, r.dbpool, query.AppendForUpdate(query.ReadEnvironment), pgx.NamedArgs{
		"envID":     envID,
		"projectID": projectID,
	})
	if err != nil {
		return svcmodel.Environment{}, fmt.Errorf("failed to read environment: %w", err)
	}

	env, err = fn(env)
	if err != nil {
		return svcmodel.Environment{}, err
	}

	_, err = r.dbpool.Exec(ctx, query.UpdateEnvironment, pgx.NamedArgs{
		"envID":      env.ID,
		"name":       env.Name,
		"serviceURL": env.ServiceURL.String(),
		"updatedAt":  env.UpdatedAt,
	})
	if err != nil {
		return svcmodel.Environment{}, fmt.Errorf("failed to update environment: %w", err)
	}

	return env, nil
}

func (r *ProjectRepository) CreateInvitation(ctx context.Context, i svcmodel.ProjectInvitation) error {
	if _, err := r.dbpool.Exec(ctx, query.CreateInvitation, pgx.NamedArgs{
		"invitationID": i.ID,
		"projectID":    i.ProjectID,
		"email":        i.Email,
		"projectRole":  i.ProjectRole,
		"tokenHash":    i.TokenHash.ToBase64(),
		"status":       i.Status,
		"createdAt":    i.CreatedAt,
		"updatedAt":    i.UpdatedAt,
	}); err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) AcceptPendingInvitation(
	ctx context.Context,
	invitationID uuid.UUID,
	fn svcmodel.AcceptProjectInvitationFunc,
) (err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	invitation, err := r.readPendingInvitationForUpdate(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("failed to read project invitation: %w", err)
	}

	// Accept the invitation
	fn(&invitation)

	_, err = tx.Exec(ctx, query.UpdateInvitation, pgx.NamedArgs{
		"invitationID": invitation.ID,
		"status":       invitation.Status,
		"updatedAt":    invitation.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to update project invitation: %w", err)
	}

	return nil
}

func (r *ProjectRepository) ReadInvitationByEmail(ctx context.Context, email string, projectID uuid.UUID) (svcmodel.ProjectInvitation, error) {
	// fetches the invitation by email for the project
	return r.readInvitation(ctx, r.dbpool, query.ReadInvitationByEmail, pgx.NamedArgs{
		"email":     email,
		"projectID": projectID,
	})
}

func (r *ProjectRepository) ReadPendingInvitationByHash(ctx context.Context, hash crypto.Hash) (svcmodel.ProjectInvitation, error) {
	return r.readInvitation(ctx, r.dbpool, query.ReadInvitationByHashAndStatus, pgx.NamedArgs{
		"hash":   hash.ToBase64(),
		"status": string(svcmodel.InvitationStatusPending),
	})
}

func (r *ProjectRepository) readPendingInvitationForUpdate(ctx context.Context, invitationID uuid.UUID) (svcmodel.ProjectInvitation, error) {
	return r.readInvitation(ctx, r.dbpool, query.ReadInvitationByIDAndStatusForUpdate, pgx.NamedArgs{
		"id":     invitationID,
		"status": string(svcmodel.InvitationStatusPending),
	})
}

func (r *ProjectRepository) readInvitation(ctx context.Context, q pgxscan.Querier, readQuery string, args pgx.NamedArgs) (svcmodel.ProjectInvitation, error) {
	var i model.ProjectInvitation

	err := pgxscan.Get(ctx, q, &i, readQuery, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.ProjectInvitation{}, apierrors.NewProjectInvitationNotFoundError().Wrap(err)
		}

		return svcmodel.ProjectInvitation{}, err
	}

	return model.ToSvcProjectInvitation(i), nil
}

func (r *ProjectRepository) ListInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error) {
	var i []model.ProjectInvitation

	err := pgxscan.Select(ctx, r.dbpool, &i, query.ListInvitationsForProject, pgx.NamedArgs{"projectID": projectID})
	if err != nil {
		return nil, err
	}

	return model.ToSvcProjectInvitations(i), nil
}

func (r *ProjectRepository) DeleteInvitation(ctx context.Context, projectID, invitationID uuid.UUID) error {
	return r.deleteInvitation(ctx, query.DeleteInvitation, pgx.NamedArgs{
		"projectID":    projectID,
		"invitationID": invitationID,
	})
}

func (r *ProjectRepository) DeleteInvitationByTokenHashAndStatus(
	ctx context.Context,
	hash crypto.Hash,
	status svcmodel.ProjectInvitationStatus,
) error {
	return r.deleteInvitation(ctx, query.DeleteInvitationByHashAndStatus, pgx.NamedArgs{
		"hash":   hash.ToBase64(),
		"status": string(status),
	})
}

func (r *ProjectRepository) deleteInvitation(ctx context.Context, deleteQuery string, args pgx.NamedArgs) error {
	result, err := r.dbpool.Exec(ctx, deleteQuery, args)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apierrors.NewProjectInvitationNotFoundError()
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

	if _, err = tx.Exec(ctx, query.CreateMember, pgx.NamedArgs{
		"userID":      m.User.ID,
		"projectID":   m.ProjectID,
		"projectRole": m.ProjectRole,
		"createdAt":   m.CreatedAt,
		"updatedAt":   m.UpdatedAt,
	}); err != nil {
		return fmt.Errorf("failed to create project member: %w", err)
	}

	if err = r.deleteInvitation(ctx, query.DeleteInvitationByEmailAndProjectID, pgx.NamedArgs{
		"email":     m.User.Email,
		"projectID": m.ProjectID,
	}); err != nil {
		return fmt.Errorf("failed to delete project invitation: %w", err)
	}

	return nil
}

func (r *ProjectRepository) ListMembersForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error) {
	var m []svcmodel.ProjectMember

	rows, err := r.dbpool.Query(ctx, query.ListMembersForProject, pgx.NamedArgs{"projectID": projectID})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		member, err := model.ScanToSvcProjectMember(rows)
		if err != nil {
			return nil, err
		}
		m = append(m, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func (r *ProjectRepository) ReadMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (svcmodel.ProjectMember, error) {
	return r.readMember(ctx, r.dbpool, query.ReadMember, pgx.NamedArgs{
		"projectID": projectID,
		"userID":    userID,
	})
}

func (r *ProjectRepository) ReadMemberByEmail(ctx context.Context, projectID uuid.UUID, email string) (svcmodel.ProjectMember, error) {
	return r.readMember(ctx, r.dbpool, query.ReadMemberByEmail, pgx.NamedArgs{
		"projectID": projectID,
		"email":     email,
	})
}

func (r *ProjectRepository) readMember(ctx context.Context, q querier, readQuery string, args pgx.NamedArgs) (svcmodel.ProjectMember, error) {
	row := q.QueryRow(ctx, readQuery, args)
	m, err := model.ScanToSvcProjectMember(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return svcmodel.ProjectMember{}, apierrors.NewProjectMemberNotFoundError().Wrap(err)
		}

		return svcmodel.ProjectMember{}, err
	}

	return m, nil
}

func (r *ProjectRepository) DeleteMember(ctx context.Context, projectID, userID uuid.UUID) error {
	result, err := r.dbpool.Exec(ctx, query.DeleteMember, pgx.NamedArgs{
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

func (r *ProjectRepository) UpdateMemberRole(
	ctx context.Context,
	projectID,
	userID uuid.UUID,
	fn svcmodel.UpdateProjectMemberFunc,
) (m svcmodel.ProjectMember, err error) {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return svcmodel.ProjectMember{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = util.FinishTransaction(ctx, tx, err)
	}()

	m, err = r.readMember(ctx, tx, query.AppendForUpdate(query.ReadMember), pgx.NamedArgs{
		"projectID": projectID,
		"userID":    userID,
	})
	if err != nil {
		return svcmodel.ProjectMember{}, fmt.Errorf("failed to read project member: %w", err)
	}

	// Update member's role
	m, err = fn(m)
	if err != nil {
		return svcmodel.ProjectMember{}, err
	}

	_, err = tx.Exec(ctx, query.UpdateMember, pgx.NamedArgs{
		"projectID":   m.ProjectID,
		"userID":      m.User.ID,
		"projectRole": m.ProjectRole,
		"updatedAt":   m.UpdatedAt,
	})
	if err != nil {
		return svcmodel.ProjectMember{}, fmt.Errorf("failed to update project member: %w", err)
	}

	return m, err
}
