package transport

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	"release-manager/transport/model"
)

func (h *Handler) listProjectInvitations(w http.ResponseWriter, r *http.Request) {
	i, err := h.ProjectInvitationSvc.ListAll(r.Context(), ContextProject(r).ID)
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetProjectInvitations(i))
}

func (h *Handler) deleteProjectInvitation(w http.ResponseWriter, r *http.Request) {
	invitationID, err := GetUUIDParamFrom(r, "invitationId")
	if err != nil {
		WriteNotFoundResponse(w, err)
		return
	}

	projectID := ContextProject(r).ID
	if _, err = h.ProjectInvitationSvc.Get(r.Context(), projectID, invitationID); err != nil {
		switch {
		case errors.Is(err, reperr.ErrResourceNotFound):
			WriteNotFoundResponse(w, err)
			return
		default:
			WriteServerErrorResponse(w, err)
			return
		}
	}

	if err := h.ProjectInvitationSvc.Delete(r.Context(), projectID, invitationID); err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
