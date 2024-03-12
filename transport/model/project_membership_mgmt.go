package model

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectMembershipManagementService interface {
	Create(ctx context.Context, r svcmodel.ProjectMembershipRequest) (svcmodel.ProjectMembershipResponse, error)
}

type ProjectMembershipRequest struct {
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

type ProjectMembershipResponse struct {
	Status   string `json:"status"`
	Resource any    `json:"resource"`
}

func ToSvcProjectMembershipRequest(r ProjectMembershipRequest, projectID, userID uuid.UUID) (svcmodel.ProjectMembershipRequest, error) {
	role, err := svcmodel.NewProjectRole(r.Role)
	if err != nil {
		return svcmodel.ProjectMembershipRequest{}, err
	}

	return svcmodel.ProjectMembershipRequest{
		ProjectID:         projectID,
		Email:             r.Email,
		Role:              role,
		RequestedByUserID: userID,
	}, nil
}

func ToNetProjectMembershipResponse(svc svcmodel.ProjectMembershipResponse) ProjectMembershipResponse {
	var r ProjectMembershipResponse
	r.Status = svc.Status

	switch t := svc.Resource.(type) {
	case svcmodel.ProjectInvitation:
		r.Resource = ToNetProjectInvitation(t)
	case svcmodel.ProjectMember:
		r.Resource = ToNetProjectMember(t)
	default:
		r.Resource = t
	}

	return r
}
