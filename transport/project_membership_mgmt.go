package transport

import (
	"errors"
	"net/http"

	svcerr "release-manager/service/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func (h *Handler) createProjectMembershipRequest(w http.ResponseWriter, r *http.Request) {
	var input model.ProjectMembershipRequest
	if err := UnmarshalRequest(r, &input); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		WriteUnprocessableEntityResponse(w, err)
		return
	}

	svcRequest, err := model.ToSvcProjectMembershipRequest(input, ContextProject(r).ID, ContextUser(r).ID)
	if err != nil {
		WriteUnprocessableEntityResponse(w, err)
	}

	resp, err := h.ProjectMembershipMgmtSvc.Create(r.Context(), svcRequest)
	if err != nil {
		switch {
		case errors.Is(err, svcerr.ErrInvitationAlreadyExists), errors.Is(err, svcerr.ErrUserIsAlreadyMember):
			WriteUnprocessableEntityResponse(w, err)
			return
		}

		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusCreated, model.ToNetProjectMembershipResponse(resp))
}
