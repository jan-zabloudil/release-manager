package handler

import (
	"net/http"

	svcmodel "release-manager/service/model"
	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) listMembers(w http.ResponseWriter, r *http.Request) {
	m, err := h.ProjectSvc.ListMembers(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
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
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateMemberRole(w http.ResponseWriter, r *http.Request) {
	var input model.UpdateProjectMemberRoleInput
	if err := util.UnmarshalRequest(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	m, err := h.ProjectSvc.UpdateMemberRole(
		r.Context(),
		svcmodel.ProjectRole(input.ProjectRole),
		util.ContextProjectID(r),
		util.ContextUserID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjectMember(m))
}
