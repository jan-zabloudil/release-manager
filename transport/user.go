package transport

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	"release-manager/transport/model"

	"github.com/google/uuid"
)

func (h *Handler) handleUser(w http.ResponseWriter, r *http.Request) {
	id, err := GetIDFromURL(r)
	if err != nil {
		WriteNotFoundResponse(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	}
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	u, err := h.UserSvc.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, reperr.ErrResourceNotFound):
			WriteNotFoundResponse(w, err)
			return
		default:
			WriteServerErrorResponse(w, err)
			return
		}
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetUser(
		u.ID,
		u.Role.Role(),
		u.Email,
		u.Name,
		u.AvatarURL,
		u.CreatedAt,
		u.UpdatedAt,
	))
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserSvc.GetAll(r.Context())
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetUsers(users))
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	if err := h.UserSvc.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, reperr.ErrResourceNotFound):
			WriteNotFoundResponse(w, err)
			return
		default:
			WriteServerErrorResponse(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
