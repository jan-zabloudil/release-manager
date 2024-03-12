package repository

import (
	"context"

	"release-manager/repository/model"
	"release-manager/repository/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type ProjectMemberRepository struct {
	client     *supabase.Client
	entity     string
	getFunc    string
	getAllFunc string
}

func NewProjectMemberRepository(c *supabase.Client) *ProjectMemberRepository {
	return &ProjectMemberRepository{
		client:     c,
		entity:     "projects_members",
		getFunc:    "get_project_member",
		getAllFunc: "get_project_members",
	}
}

func (r *ProjectMemberRepository) Read(ctx context.Context, projectID, userID uuid.UUID) (svcmodel.ProjectMember, error) {
	input := map[string]interface{}{
		"p_project_id": projectID,
		"p_user_id":    userID,
	}

	var resp []model.ProjectMember
	err := r.client.
		DB.Rpc(r.getFunc, input).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return svcmodel.ProjectMember{}, utils.WrapSupabaseDBErr(err)
	}

	if err := utils.ValidateSingleRecordFetchAfterReadOperation(resp); err != nil {
		return svcmodel.ProjectMember{}, err
	}

	return model.ToSvcProjectMember(resp[0])
}

func (r *ProjectMemberRepository) ReadAll(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error) {
	input := map[string]interface{}{
		"p_project_id": projectID,
	}

	var resp []model.ProjectMember
	err := r.client.
		DB.Rpc(r.getAllFunc, input).
		ExecuteWithContext(ctx, &resp)
	if err != nil {
		return nil, utils.WrapSupabaseDBErr(err)
	}

	return model.ToSvcProjectMembers(resp)
}

func (r *ProjectMemberRepository) Insert(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, role svcmodel.ProjectRole, invitedByUserID uuid.UUID) (svcmodel.ProjectMember, error) {
	err := r.client.DB.
		From(r.entity).
		Insert(model.ToProjectMemberInput(projectID, userID, role, invitedByUserID)).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return svcmodel.ProjectMember{}, utils.WrapSupabaseDBErr(err)
	}

	return r.Read(ctx, projectID, userID)
}

func (r *ProjectMemberRepository) Delete(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {

	err := r.client.DB.
		From(r.entity).
		Delete().
		Eq("project_id", projectID.String()).
		Eq("user_id", userID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return utils.WrapSupabaseDBErr(err)
	}

	return nil
}

func (r *ProjectMemberRepository) Update(ctx context.Context, m svcmodel.ProjectMember) (svcmodel.ProjectMember, error) {

	err := r.client.
		DB.From(r.entity).
		Update(model.ProjectInvitationPatch{Role: m.Role.Role()}).
		Eq("project_id", m.ProjectID.String()).
		Eq("user_id", m.User.ID.String()).
		ExecuteWithContext(ctx, nil)
	if err != nil {
		return svcmodel.ProjectMember{}, utils.WrapSupabaseDBErr(err)
	}

	return r.Read(ctx, m.ProjectID, m.User.ID)
}
