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
	id, err := GetUUIDParamFrom(r, "id")
	if err != nil {
		WriteNotFoundResponse(w, err)
		return
	}

	if _, err = h.ProjectInvitationSvc.Get(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, reperr.ErrResourceNotFound):
			WriteNotFoundResponse(w, err)
			return
		default:
			WriteServerErrorResponse(w, err)
			return
		}
	}

	if err := h.ProjectInvitationSvc.Delete(r.Context(), id); err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
