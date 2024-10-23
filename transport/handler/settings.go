package handler

import (
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateSettingsInput
	if err := util.UnmarshalBody(r, &req); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.SettingsSvc.Update(
		r.Context(),
		model.ToSvcUpdateSettingsInput(req),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.SettingsSvc.Get(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToSettings(s))
}
