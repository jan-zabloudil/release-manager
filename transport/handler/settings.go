package handler

import (
	"net/http"

	resperr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var input model.UpdateSettingsInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewFromBodyUnmarshalErr(err))
		return
	}

	if err := h.SettingsSvc.Update(
		r.Context(),
		model.ToSvcUpdateSettingsInput(input),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.SettingsSvc.Get(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToSettings(s))
}
