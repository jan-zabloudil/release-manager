package handler

import (
	"errors"
	"net/http"

	svcerr "release-manager/service/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func (h *Handler) createProjectMembershipRequest(w http.ResponseWriter, r *http.Request) {
	var input model.ProjectMembershipRequest
	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	svcRequest, err := model.ToSvcProjectMembershipRequest(input, utils.ContextProject(r).ID, utils.ContextUser(r).ID)
	if err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
	}

	resp, err := h.ProjectMembershipMgmtSvc.Create(r.Context(), svcRequest, utils.ContextProjectMember(r))
	if err != nil {
		switch {
		case errors.Is(err, svcerr.ErrInvitationAlreadyExists), errors.Is(err, svcerr.ErrUserIsAlreadyMember):
			utils.WriteUnprocessableEntityResponse(w, err)
			return
		case errors.Is(err, svcerr.ErrProjectMemberRoleCannotBeGranted):
			utils.WriteForbiddenErrorResponse(w, err)
			return
		}

		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, model.ToNetProjectMembershipResponse(resp))
}
