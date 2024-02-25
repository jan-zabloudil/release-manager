package transport

import (
	"net/http"

	"release-manager/transport/model"
)

func (h *Handler) setSettings(w http.ResponseWriter, r *http.Request) {
	var input model.Settings
	if err := UnmarshalRequest(r, &input); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	s, err := h.SettingsSvc.Get(r.Context())
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	s, err = h.SettingsSvc.Set(r.Context(), model.ToSvcSettings(
		s,
		input.OrganizationName,
		input.SlackToken,
		input.GithubToken,
		input.DefaultReleaseMsg,
	))
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetSettings(
		s.OrganizationName,
		s.SlackToken,
		s.GithubToken,
		s.DefaultReleaseMsg,
	))
}

func (h *Handler) getSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.SettingsSvc.Get(r.Context())
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetSettings(
		s.OrganizationName,
		s.SlackToken,
		s.GithubToken,
		s.DefaultReleaseMsg,
	))
}
