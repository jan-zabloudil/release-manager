package handler

import (
	"net/http"

	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) listMembers(w http.ResponseWriter, r *http.Request) {
	m, err := h.ProjectSvc.ListMembers(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjectMembers(m))
}

func (h *Handler) deleteMember(w http.ResponseWriter, r *http.Request) {
	err := h.ProjectSvc.DeleteMember(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextUserID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
