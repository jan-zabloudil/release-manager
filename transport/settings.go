package transport

import (
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateSettingsInput
	if err := UnmarshalRequest(r, &req); err != nil {
		WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	s, err := h.SettingsSvc.Update(
		r.Context(),
		model.ToSvcUpdateSettingsInput(req),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToSettings(s))
}

func (h *Handler) getSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.SettingsSvc.Get(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToSettings(s))
}
