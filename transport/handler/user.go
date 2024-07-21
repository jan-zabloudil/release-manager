package handler

import (
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) getAuthUser(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserSvc.Get(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToUser(u))
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserSvc.GetForAdmin(
		r.Context(),
		util.ContextUserID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToUser(u))
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserSvc.ListAllForAdmin(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToUsers(users))
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.UserSvc.DeleteForAdmin(
		r.Context(),
		util.ContextUserID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
