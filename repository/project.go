package repository

import (
	"context"
	"fmt"

	"release-manager/pkg/crypto"
	"release-manager/pkg/id"
	"release-manager/repository/helper"
	"release-manager/repository/model"
	"release-manager/repository/query"
	svcerrors "release-manager/service/errors"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueEnvironmentNamePerProjectConstraintName = "unique_environment_name_per_project"
	uniqueInvitationPerProjectConstraintName      = "unique_invitation_per_project"
	uniqueGithubRepoConstraintName                = "unique_github_repo"
)

type ProjectRepository struct {
	dbpool             *pgxpool.Pool
	githubURLGenerator githubURLGenerator
}

func NewProjectRepository(pool *pgxpool.Pool, urlGenerator githubURLGenerator) *ProjectRepository {
	return &ProjectRepository{
		dbpool:             pool,
		githubURLGenerator: urlGenerator,
	}
}

func (r *ProjectRepository) CreateProjectWithOwner(ctx context.Context, p svcmodel.Project, owner svcmodel.ProjectMember) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, query.CreateProject, pgx.NamedArgs{
			"id":             p.ID,
			"name":           p.Name,
			"slackChannelID": p.SlackChannelID,
			// convert to db model in order to correctly save the struct to json field
			"releaseNotificationConfig": model.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
			"createdAt":                 p.CreatedAt,
			"updatedAt":                 p.UpdatedAt,
		}); err != nil {
			return fmt.Errorf("creating project: %w", err)
		}

		if _, err := tx.Exec(ctx, query.CreateMember, pgx.NamedArgs{
			"userID":      owner.User.ID,
			"projectID":   p.ID,
			"projectRole": owner.ProjectRole,
			"createdAt":   owner.CreatedAt,
			"updatedAt":   owner.UpdatedAt,
		}); err != nil {
			return fmt.Errorf("creating project member: %w", err)
		}

		return nil
	})
}

func (r *ProjectRepository) ReadProject(ctx context.Context, id uuid.UUID) (svcmodel.Project, error) {
	return r.readProject(ctx, r.dbpool, query.ReadProject, pgx.NamedArgs{"id": id})
}

func (r *ProjectRepository) ListProjects(ctx context.Context) ([]svcmodel.Project, error) {
	return r.listProjects(ctx, query.ListProjects, nil)
}

func (r *ProjectRepository) ListProjectsForUser(ctx context.Context, userID id.AuthUser) ([]svcmodel.Project, error) {
	return r.listProjects(ctx, query.ListProjectsForUser, pgx.NamedArgs{"userID": userID})
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	result, err := r.dbpool.Exec(ctx, query.DeleteProject, pgx.NamedArgs{"id": id})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return svcerrors.NewProjectNotFoundError()
	}

	return nil
}

func (r *ProjectRepository) UpdateProject(
	ctx context.Context,
	projectID uuid.UUID,
	updateFn func(p svcmodel.Project) (svcmodel.Project, error),
) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		p, err := r.readProject(ctx, tx, query.AppendForUpdate(query.ReadProject), pgx.NamedArgs{"id": projectID})
		if err != nil {
			return fmt.Errorf("reading project: %w", err)
		}

		p, err = updateFn(p)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, query.UpdateProject, pgx.NamedArgs{
			"id":             p.ID,
			"name":           p.Name,
			"slackChannelID": p.SlackChannelID,
			// convert to db model in order to correctly save the struct to json field
			"releaseNotificationConfig": model.ReleaseNotificationConfig(p.ReleaseNotificationConfig),
			"githubOwnerSlug":           p.GithubOwnerSlug(),
			"githubRepoSlug":            p.GithubRepoSlug(),
			"updatedAt":                 p.UpdatedAt,
		}); err != nil {
			if helper.IsUniqueConstraintViolation(err, uniqueGithubRepoConstraintName) {
				return svcerrors.NewProjectGithubRepoAlreadyUsedError().Wrap(err)
			}

			return fmt.Errorf("updating project: %w", err)
		}

		return nil
	})
}

func (r *ProjectRepository) CreateEnvironment(ctx context.Context, e svcmodel.Environment) error {
	if _, err := r.dbpool.Exec(ctx, query.CreateEnvironment, pgx.NamedArgs{
		"id":         e.ID,
		"projectID":  e.ProjectID,
		"name":       e.Name,
		"serviceURL": e.ServiceURL.String(),
		"createdAt":  e.CreatedAt,
		"updatedAt":  e.UpdatedAt,
	}); err != nil {
		if helper.IsUniqueConstraintViolation(err, uniqueEnvironmentNamePerProjectConstraintName) {
			return svcerrors.NewEnvironmentDuplicateNameError().Wrap(err)
		}

		return err
	}

	return nil
}

func (r *ProjectRepository) ReadEnvironment(ctx context.Context, projectID uuid.UUID, envID id.Environment) (svcmodel.Environment, error) {
	return r.readEnvironment(ctx, r.dbpool, query.ReadEnvironment, pgx.NamedArgs{
		"projectID": projectID,
		"envID":     envID,
	})
}

func (r *ProjectRepository) ListEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Environment, error) {
	e, err := helper.ListValues[model.Environment](ctx, r.dbpool, query.ListEnvironmentsForProject, pgx.NamedArgs{
		"projectID": projectID},
	)
	if err != nil {
		return nil, err
	}

	return model.ToSvcEnvironments(e)
}

func (r *ProjectRepository) DeleteEnvironment(ctx context.Context, projectID uuid.UUID, envID id.Environment) error {
	result, err := r.dbpool.Exec(ctx, query.DeleteEnvironment, pgx.NamedArgs{
		"envID":     envID,
		"projectID": projectID,
	})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return svcerrors.NewEnvironmentNotFoundError()
	}

	return nil
}

func (r *ProjectRepository) UpdateEnvironment(
	ctx context.Context,
	projectID uuid.UUID,
	envID id.Environment,
	updateFn func(e svcmodel.Environment) (svcmodel.Environment, error),
) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		env, err := r.readEnvironment(ctx, tx, query.AppendForUpdate(query.ReadEnvironment), pgx.NamedArgs{
			"projectID": projectID,
			"envID":     envID,
		})
		if err != nil {
			return fmt.Errorf("reading environment: %w", err)
		}

		env, err = updateFn(env)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(ctx, query.UpdateEnvironment, pgx.NamedArgs{
			"envID":      env.ID,
			"name":       env.Name,
			"serviceURL": env.ServiceURL.String(),
			"updatedAt":  env.UpdatedAt,
		}); err != nil {
			if helper.IsUniqueConstraintViolation(err, uniqueEnvironmentNamePerProjectConstraintName) {
				return svcerrors.NewEnvironmentDuplicateNameError().Wrap(err)
			}

			return fmt.Errorf("updating environment: %w", err)
		}

		return nil
	})
}

func (r *ProjectRepository) CreateInvitation(ctx context.Context, i svcmodel.ProjectInvitation) error {
	if _, err := r.dbpool.Exec(ctx, query.CreateInvitation, pgx.NamedArgs{
		"invitationID": i.ID,
		"projectID":    i.ProjectID,
		"email":        i.Email,
		"projectRole":  i.ProjectRole,
		"tokenHash":    i.TokenHash.ToBase64(),
		"status":       i.Status,
		"invitedBy":    i.InviterUserID,
		"createdAt":    i.CreatedAt,
		"updatedAt":    i.UpdatedAt,
	}); err != nil {
		if helper.IsUniqueConstraintViolation(err, uniqueInvitationPerProjectConstraintName) {
			return svcerrors.NewProjectInvitationAlreadyExistsError().Wrap(err)
		}

		return err
	}

	return nil
}

func (r *ProjectRepository) ListInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error) {
	i, err := helper.ListValues[model.ProjectInvitation](ctx, r.dbpool, query.ListInvitationsForProject, pgx.NamedArgs{
		"projectID": projectID},
	)
	if err != nil {
		return nil, err
	}

	return model.ToSvcProjectInvitations(i), nil
}

func (r *ProjectRepository) DeleteInvitation(ctx context.Context, projectID uuid.UUID, invitationID id.ProjectInvitation) error {
	return r.deleteInvitation(ctx, r.dbpool, query.DeleteInvitation, pgx.NamedArgs{
		"projectID":    projectID,
		"invitationID": invitationID,
	})
}

func (r *ProjectRepository) DeleteInvitationByTokenHashAndStatus(
	ctx context.Context,
	hash crypto.Hash,
	status svcmodel.ProjectInvitationStatus,
) error {
	return r.deleteInvitation(ctx, r.dbpool, query.DeleteInvitationByHashAndStatus, pgx.NamedArgs{
		"hash":   hash.ToBase64(),
		"status": string(status),
	})
}

func (r *ProjectRepository) UpdateInvitation(
	ctx context.Context,
	invitationHash crypto.Hash,
	updateFn func(i svcmodel.ProjectInvitation) (svcmodel.ProjectInvitation, error),
) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		i, err := r.readInvitation(ctx, tx, query.AppendForUpdate(query.ReadInvitationByHash), pgx.NamedArgs{
			"hash": invitationHash.ToBase64(),
		})
		if err != nil {
			return fmt.Errorf("reading project invitation: %w", err)
		}

		i, err = updateFn(i)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, query.UpdateInvitation, pgx.NamedArgs{
			"invitationID": i.ID,
			"status":       i.Status,
			"updatedAt":    i.UpdatedAt,
		}); err != nil {
			return fmt.Errorf("updating project invitation: %w", err)
		}

		return nil
	})
}

// CreateMember creates a project member and deletes the invitation
func (r *ProjectRepository) CreateMember(
	ctx context.Context,
	tokenHash crypto.Hash,
	createMemberFn func(i svcmodel.ProjectInvitation) (svcmodel.ProjectMember, error),
) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		i, err := r.readInvitation(ctx, tx, query.ReadInvitationByHash, pgx.NamedArgs{
			"hash": tokenHash.ToBase64(),
		})
		if err != nil {
			return fmt.Errorf("reading project invitation: %w", err)
		}

		m, err := createMemberFn(i)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, query.CreateMember, pgx.NamedArgs{
			"userID":      m.User.ID,
			"projectID":   m.ProjectID,
			"projectRole": m.ProjectRole,
			"createdAt":   m.CreatedAt,
			"updatedAt":   m.UpdatedAt,
		}); err != nil {
			return fmt.Errorf("creating project member: %w", err)
		}

		if err := r.deleteInvitation(ctx, tx, query.DeleteInvitation, pgx.NamedArgs{
			"projectID":    i.ProjectID,
			"invitationID": i.ID,
		}); err != nil {
			return fmt.Errorf("deleting project invitation: %w", err)
		}

		return nil
	})
}

func (r *ProjectRepository) ListMembersForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error) {
	return r.listMembers(ctx, r.dbpool, query.ListMembersForProject, pgx.NamedArgs{"projectID": projectID})
}

func (r *ProjectRepository) ListMembersForUser(ctx context.Context, userID id.AuthUser) ([]svcmodel.ProjectMember, error) {
	return r.listMembers(ctx, r.dbpool, query.ListMembersForUser, pgx.NamedArgs{"userID": userID})
}

func (r *ProjectRepository) ReadMember(ctx context.Context, projectID uuid.UUID, userID id.User) (svcmodel.ProjectMember, error) {
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

func (r *ProjectRepository) DeleteMember(ctx context.Context, projectID uuid.UUID, userID id.User) error {
	result, err := r.dbpool.Exec(ctx, query.DeleteMember, pgx.NamedArgs{
		"projectID": projectID,
		"userID":    userID,
	})
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return svcerrors.NewProjectMemberNotFoundError()
	}

	return nil
}

func (r *ProjectRepository) UpdateMember(
	ctx context.Context,
	projectID uuid.UUID,
	userID id.User,
	updateFn func(m svcmodel.ProjectMember) (svcmodel.ProjectMember, error),
) error {
	return helper.RunTransaction(ctx, r.dbpool, func(tx pgx.Tx) error {
		m, err := r.readMember(ctx, tx, query.AppendForUpdate(query.ReadMember), pgx.NamedArgs{
			"projectID": projectID,
			"userID":    userID,
		})
		if err != nil {
			return fmt.Errorf("reading project member: %w", err)
		}

		m, err = updateFn(m)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, query.UpdateMember, pgx.NamedArgs{
			"projectID":   m.ProjectID,
			"userID":      m.User.ID,
			"projectRole": m.ProjectRole,
			"updatedAt":   m.UpdatedAt,
		}); err != nil {
			return fmt.Errorf("updating project member: %w", err)
		}

		return err
	})
}

func (r *ProjectRepository) readProject(ctx context.Context, q helper.Querier, query string, args pgx.NamedArgs) (svcmodel.Project, error) {
	p, err := helper.ReadValue[model.Project](ctx, q, query, args)
	if err != nil {
		if helper.IsNotFound(err) {
			return svcmodel.Project{}, svcerrors.NewProjectNotFoundError().Wrap(err)
		}

		return svcmodel.Project{}, err
	}

	return model.ToSvcProject(p, r.githubURLGenerator.GenerateRepoURL)
}

func (r *ProjectRepository) listProjects(ctx context.Context, query string, args pgx.NamedArgs) ([]svcmodel.Project, error) {
	p, err := helper.ListValues[model.Project](ctx, r.dbpool, query, args)
	if err != nil {
		return nil, err
	}

	return model.ToSvcProjects(p, r.githubURLGenerator.GenerateRepoURL)
}

func (r *ProjectRepository) readEnvironment(ctx context.Context, q helper.Querier, query string, args pgx.NamedArgs) (svcmodel.Environment, error) {
	e, err := helper.ReadValue[model.Environment](ctx, q, query, args)
	if err != nil {
		if helper.IsNotFound(err) {
			return svcmodel.Environment{}, svcerrors.NewEnvironmentNotFoundError().Wrap(err)
		}

		return svcmodel.Environment{}, err
	}

	return model.ToSvcEnvironment(e)
}

func (r *ProjectRepository) readInvitation(ctx context.Context, q helper.Querier, query string, args pgx.NamedArgs) (svcmodel.ProjectInvitation, error) {
	i, err := helper.ReadValue[model.ProjectInvitation](ctx, q, query, args)
	if err != nil {
		if helper.IsNotFound(err) {
			return svcmodel.ProjectInvitation{}, svcerrors.NewProjectInvitationNotFoundError().Wrap(err)
		}

		return svcmodel.ProjectInvitation{}, err
	}

	return model.ToSvcProjectInvitation(i), nil
}

func (r *ProjectRepository) deleteInvitation(ctx context.Context, e helper.ExecExecutor, query string, args pgx.NamedArgs) error {
	result, err := e.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return svcerrors.NewProjectInvitationNotFoundError()
	}

	return nil
}

func (r *ProjectRepository) listMembers(ctx context.Context, q helper.Querier, query string, args pgx.NamedArgs) ([]svcmodel.ProjectMember, error) {
	m, err := helper.ListValues[model.ProjectMember](ctx, q, query, args)
	if err != nil {
		return nil, err
	}

	return model.ToSvcProjectMembers(m), nil
}

func (r *ProjectRepository) readMember(ctx context.Context, q helper.Querier, query string, args pgx.NamedArgs) (svcmodel.ProjectMember, error) {
	m, err := helper.ReadValue[model.ProjectMember](ctx, q, query, args)
	if err != nil {
		if helper.IsNotFound(err) {
			return svcmodel.ProjectMember{}, svcerrors.NewProjectMemberNotFoundError().Wrap(err)
		}

		return svcmodel.ProjectMember{}, err
	}

	return model.ToSvcProjectMember(m), nil
}
