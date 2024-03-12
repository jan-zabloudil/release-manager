package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type ProjectInvitationRepository struct {
	client     *supabase.Client
	entity     string
	insertFunc string
}

func NewProjectInvitationRepository(c *supabase.Client) *ProjectInvitationRepository {
	return &ProjectInvitationRepository{
		client:     c,
		entity:     "projects_invitations",
		insertFunc: "create_project_invitation",
	}
}

func (r *ProjectInvitationRepository) Insert(ctx context.Context, projectID uuid.UUID, email string, role svcmodel.ProjectRole, invitedByUserID uuid.UUID) (svcmodel.ProjectInvitation, error) {

	var resp []model.ProjectInvitationResponse
	err := r.client.
		DB.Rpc(r.insertFunc, model.ToProjectInvitationInput(projectID, email, role, invitedByUserID)).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectInvitation{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetchAfterWriteOperation(resp); err != nil {
		return svcmodel.ProjectInvitation{}, err
	}

	return model.ToSvcProjectInvitation(resp[0])
}

func (r *ProjectInvitationRepository) ReadByEmail(ctx context.Context, projectID uuid.UUID, email string) (svcmodel.ProjectInvitation, error) {

	var resp model.ProjectInvitationResponse
	err := r.client.DB.
		From(r.entity).
		Select("*").Single().
		Eq("project_id", projectID.String()).
		Eq("email", email).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectInvitation{}, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcProjectInvitation(resp)
}

func (r *ProjectInvitationRepository) Read(ctx context.Context, projectID, invitationID uuid.UUID) (svcmodel.ProjectInvitation, error) {

	var resp model.ProjectInvitationResponse
	err := r.client.DB.
		From(r.entity).
		Select("*").Single().
		Eq("id", invitationID.String()).
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectInvitation{}, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcProjectInvitation(resp)
}

func (r *ProjectInvitationRepository) ReadAll(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error) {

	var resp []model.ProjectInvitationResponse
	err := r.client.DB.
		From(r.entity).
		Select("*").
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcProjectInvitations(resp)
}

func (r *ProjectInvitationRepository) Delete(ctx context.Context, projectID, invitationID uuid.UUID) error {

	err := r.client.DB.
		From(r.entity).
		Delete().
		Eq("id", invitationID.String()).
		Eq("project_id", projectID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return utils.WrapSupabaseDBErr(err)
	}

	return nil
}
