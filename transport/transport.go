package transport

import (
	"log/slog"
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"

	"github.com/go-chi/chi/v5"
	httpx "go.strv.io/net/http"
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
	h := Handler{
		Mux:                  chi.NewRouter(),
		AuthSvc:              as,
		UserSvc:              us,
		ProjectSvc:           ps,
		SettingsSvc:          ss,
		ProjectMembershipSvc: pms,
	}

	h.Mux.Use(httpx.RequestIDMiddleware(RequestID))
	h.Mux.Use(httpx.LoggingMiddleware(slog.Default().WithGroup("logger")))
	h.Mux.Use(httpx.RecoverMiddleware(slog.Default().WithGroup("recover")))
	h.Mux.Use(h.auth)

	h.Mux.Get("/admin/users", h.requireAuthUser(h.listUsers))
	h.Mux.Route("/admin/users/{id}", func(r chi.Router) {
		r.Use(h.handleResourceID("id", util.ContextSetUserID))
		r.Get("/", h.requireAuthUser(h.getUser))
		r.Delete("/", h.requireAuthUser(h.deleteUser))
	})

	h.Mux.Route("/projects", func(r chi.Router) {
		r.Post("/", h.requireAuthUser(h.createProject))
		r.Get("/", h.requireAuthUser(h.listProjects))
		r.Route("/{id}", func(r chi.Router) {
			r.Use(h.handleResourceID("id", util.ContextSetProjectID))
			r.Get("/", h.requireAuthUser(h.getProject))
			r.Patch("/", h.requireAuthUser(h.updateProject))
			r.Delete("/", h.requireAuthUser(h.deleteProject))
			r.Route("/environments", func(r chi.Router) {
				r.Post("/", h.requireAuthUser(h.createEnvironment))
				r.Get("/", h.requireAuthUser(h.listEnvironments))
				r.Route("/{environment_id}", func(r chi.Router) {
					r.Use(h.handleResourceID("environment_id", util.ContextSetEnvironmentID))
					r.Get("/", h.requireAuthUser(h.getEnvironment))
					r.Patch("/", h.requireAuthUser(h.updateEnvironment))
					r.Delete("/", h.requireAuthUser(h.deleteEnvironment))
				})
			})
			r.Route("/invitations", func(r chi.Router) {
				r.Post("/", h.requireAuthUser(h.createInvitation))
				r.Get("/", h.requireAuthUser(h.listInvitations))
				r.Route("/{invitation_id}", func(r chi.Router) {
					r.Use(h.handleResourceID("invitation_id", util.ContextSetProjectInvitationID))
					r.Delete("/", h.requireAuthUser(h.deleteInvitation))
				})
			})
		})
	})

	h.Mux.Route("/projects/invitations", func(r chi.Router) {
		r.Use(h.handleInvitationToken)
		r.Get("/accept", h.acceptInvitation)
		r.Get("/reject", h.rejectInvitation)
	})

	h.Mux.Route("/organization/settings", func(r chi.Router) {
		r.Get("/", h.requireAuthUser(h.getSettings))
		r.Patch("/", h.requireAuthUser(h.updateSettings))
	})

	h.Mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.Mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		WriteResponseError(w, responseerrors.NewNotFoundError())
	})
	h.Mux.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		WriteResponseError(w, responseerrors.NewMethodNotAllowedError())
	})

	return &h
}
