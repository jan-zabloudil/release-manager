package transport

import (
	"net/http"

	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserSvc.Get(
		r.Context(),
		util.ContextUserID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToUser(
		u.ID,
		u.Role,
		u.Email,
		u.Name,
		u.AvatarURL,
		u.CreatedAt,
		u.UpdatedAt,
	))
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserSvc.GetAll(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToUsers(users))
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.UserSvc.Delete(
		r.Context(),
		util.ContextUserID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
