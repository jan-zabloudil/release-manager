package handler

import (
	"release-manager/transport/model"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Mux                  *chi.Mux
	AuthSvc              model.AuthService
	UserSvc              model.UserService
	ProjectSvc           model.ProjectService
	SettingsSvc          model.SettingsService
	ProjectMembershipSvc model.ProjectMembershipService
}

func NewHandler(
	as model.AuthService,
	us model.UserService,
	ps model.ProjectService,
	ss model.SettingsService,
	pms model.ProjectMembershipService,
) *Handler {
	h := &Handler{
		Mux:                  chi.NewRouter(),
		AuthSvc:              as,
		UserSvc:              us,
		ProjectSvc:           ps,
		SettingsSvc:          ss,
		ProjectMembershipSvc: pms,
	}

	h.setupRoutes()

	return h
}
