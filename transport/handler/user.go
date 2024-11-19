package handler

import (
	"net/http"

	"release-manager/pkg/id"
	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) getAuthUser(w http.ResponseWriter, r *http.Request) {
	authUserID := util.ContextAuthUserID(r)

	u, err := h.UserSvc.GetAuthenticated(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	m, err := h.ProjectSvc.ListMembersForUser(r.Context(), authUserID)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToAuthUser(u, m))
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetPathParam[id.User](r, "id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	u, err := h.UserSvc.GetForAdmin(
		r.Context(),
		userID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToUser(u))
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserSvc.ListAllForAdmin(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToUsers(users))
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetPathParam[id.User](r, "id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.UserSvc.DeleteForAdmin(
		r.Context(),
		userID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
