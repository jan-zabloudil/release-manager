package handler

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func (h *Handler) listProjectInvitations(w http.ResponseWriter, r *http.Request) {
	i, err := h.ProjectInvitationSvc.ListAll(r.Context(), utils.ContextProject(r).ID)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetProjectInvitations(i))
}

func (h *Handler) deleteProjectInvitation(w http.ResponseWriter, r *http.Request) {
	invitationID, err := utils.GetUUIDParamFrom(r, "invitationId")
	if err != nil {
		utils.WriteNotFoundResponse(w, err)
		return
	}

	projectID := utils.ContextProject(r).ID
	if _, err = h.ProjectInvitationSvc.Get(r.Context(), projectID, invitationID); err != nil {
		switch {
		case errors.Is(err, reperr.ErrResourceNotFound):
			utils.WriteNotFoundResponse(w, err)
			return
		default:
			utils.WriteServerErrorResponse(w, err)
			return
		}
	}

	if err := h.ProjectInvitationSvc.Delete(r.Context(), projectID, invitationID); err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
