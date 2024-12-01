package handler

import (
	"net/http"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
	resperr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createInvitation(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.CreateProjectInvitationInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewFromBodyUnmarshalErr(err))
		return
	}

	i, err := h.ProjectSvc.Invite(
		r.Context(),
		model.ToSvcCreateProjectInvitationInput(input, projectID),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToProjectInvitation(i))
}

func (h *Handler) listInvitations(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	i, err := h.ProjectSvc.ListInvitations(r.Context(), projectID, util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjectInvitations(i))
}

func (h *Handler) cancelInvitation(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.CancelProjectInvitationURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromURLParamsUnmarshalErr(err))
		return
	}

	if err := h.ProjectSvc.CancelInvitation(
		r.Context(),
		params.ProjectID,
		params.InvitationID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) acceptInvitation(w http.ResponseWriter, r *http.Request) {
	param, err := util.UnmarshalURLParams[model.ProjectInvitationTokenURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromURLParamsUnmarshalErr(err))
		return
	}

	if err := h.ProjectSvc.AcceptInvitation(
		r.Context(),
		svcmodel.ProjectInvitationToken(param.Token),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) rejectInvitation(w http.ResponseWriter, r *http.Request) {
	param, err := util.UnmarshalURLParams[model.ProjectInvitationTokenURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromURLParamsUnmarshalErr(err))
		return
	}

	if err := h.ProjectSvc.RejectInvitation(
		r.Context(),
		svcmodel.ProjectInvitationToken(param.Token),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
