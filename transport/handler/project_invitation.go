package handler

import (
	"net/http"

	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createInvitation(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.CreateProjectInvitationInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	i, err := h.ProjectSvc.Invite(
		r.Context(),
		model.ToSvcCreateProjectInvitationInput(input, projectID),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToProjectInvitation(i))
}

func (h *Handler) listInvitations(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	i, err := h.ProjectSvc.ListInvitations(r.Context(), projectID, util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjectInvitations(i))
}

func (h *Handler) cancelInvitation(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.ProjectInvitationURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.CancelInvitation(
		r.Context(),
		params.ProjectID,
		params.InvitationID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) acceptInvitation(w http.ResponseWriter, r *http.Request) {
	tkn, err := util.GetQueryParam[cryptox.Token](r, "token")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.AcceptInvitation(
		r.Context(),
		svcmodel.ProjectInvitationToken(tkn),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) rejectInvitation(w http.ResponseWriter, r *http.Request) {
	tkn, err := util.GetQueryParam[cryptox.Token](r, "token")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.RejectInvitation(
		r.Context(),
		svcmodel.ProjectInvitationToken(tkn),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
