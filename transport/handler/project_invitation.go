package handler

import (
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createInvitation(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProjectInvitationInput
	if err := util.UnmarshalRequest(r, &req); err != nil {
		util.WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	i, err := h.ProjectMembershipSvc.CreateInvitation(
		r.Context(),
		model.ToSvcCreateProjectInvitationInput(req, util.ContextProjectID(r)),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToProjectInvitation(i))
}

func (h *Handler) listInvitations(w http.ResponseWriter, r *http.Request) {
	i, err := h.ProjectMembershipSvc.ListInvitations(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjectInvitations(i))
}

func (h *Handler) deleteInvitation(w http.ResponseWriter, r *http.Request) {
	err := h.ProjectMembershipSvc.DeleteInvitation(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextProjectInvitationID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) acceptInvitation(w http.ResponseWriter, r *http.Request) {
	err := h.ProjectMembershipSvc.AcceptInvitation(r.Context(), util.ContextProjectInvitationToken(r))
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) rejectInvitation(w http.ResponseWriter, r *http.Request) {
	err := h.ProjectMembershipSvc.RejectInvitation(r.Context(), util.ContextProjectInvitationToken(r))
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
