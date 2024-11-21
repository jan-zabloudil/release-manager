package handler

import (
	"net/http"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
	resperr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) listMembers(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	m, err := h.ProjectSvc.ListMembersForProject(r.Context(), projectID, util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjectMembers(m))
}

func (h *Handler) deleteMember(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.ProjectMemberURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.DeleteMember(
		r.Context(),
		params.ProjectID,
		params.UserID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateMemberRole(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.ProjectMemberURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.UpdateProjectMemberRoleInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewFromBodyUnmarshalErr(err))
		return
	}

	if err := h.ProjectSvc.UpdateMemberRole(
		r.Context(),
		svcmodel.ProjectRole(input.ProjectRole),
		params.ProjectID,
		params.UserID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
