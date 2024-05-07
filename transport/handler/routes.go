package handler

import (
	"log/slog"
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/middleware"
	"release-manager/transport/util"

	"github.com/go-chi/chi/v5"
	httpx "go.strv.io/net/http"
)

func (h *Handler) setupRoutes() {
	h.Mux.Use(httpx.RequestIDMiddleware(util.RequestID))
	h.Mux.Use(httpx.LoggingMiddleware(slog.Default().WithGroup("logger")))
	h.Mux.Use(httpx.RecoverMiddleware(slog.Default().WithGroup("recover")))
	h.Mux.Use(middleware.Auth(h.AuthSvc))

	h.Mux.Get("/admin/users", middleware.RequireAuthUser(h.listUsers))
	h.Mux.Route("/admin/users/{id}", func(r chi.Router) {
		r.Use(middleware.HandleResourceID("id", util.ContextSetUserID))
		r.Get("/", middleware.RequireAuthUser(h.getUser))
		r.Delete("/", middleware.RequireAuthUser(h.deleteUser))
	})

	h.Mux.Route("/projects", func(r chi.Router) {
		r.Post("/", middleware.RequireAuthUser(h.createProject))
		r.Get("/", middleware.RequireAuthUser(h.listProjects))
		r.Route("/{id}", func(r chi.Router) {
			r.Use(middleware.HandleResourceID("id", util.ContextSetProjectID))
			r.Get("/", middleware.RequireAuthUser(h.getProject))
			r.Patch("/", middleware.RequireAuthUser(h.updateProject))
			r.Delete("/", middleware.RequireAuthUser(h.deleteProject))
			r.Route("/environments", func(r chi.Router) {
				r.Post("/", middleware.RequireAuthUser(h.createEnvironment))
				r.Get("/", middleware.RequireAuthUser(h.listEnvironments))
				r.Route("/{environment_id}", func(r chi.Router) {
					r.Use(middleware.HandleResourceID("environment_id", util.ContextSetEnvironmentID))
					r.Get("/", middleware.RequireAuthUser(h.getEnvironment))
					r.Patch("/", middleware.RequireAuthUser(h.updateEnvironment))
					r.Delete("/", middleware.RequireAuthUser(h.deleteEnvironment))
				})
			})
			r.Route("/invitations", func(r chi.Router) {
				r.Post("/", middleware.RequireAuthUser(h.createInvitation))
				r.Get("/", middleware.RequireAuthUser(h.listInvitations))
				r.Route("/{invitation_id}", func(r chi.Router) {
					r.Use(middleware.HandleResourceID("invitation_id", util.ContextSetProjectInvitationID))
					r.Delete("/", middleware.RequireAuthUser(h.cancelInvitation))
				})
			})

			r.Route("/members", func(r chi.Router) {
				r.Get("/", middleware.RequireAuthUser(h.listMembers))
				r.Route("/{user_id}", func(r chi.Router) {
					r.Use(middleware.HandleResourceID("user_id", util.ContextSetUserID))
					r.Delete("/", middleware.RequireAuthUser(h.deleteMember))
				})
			})

			r.Get("/repository/tags", middleware.RequireAuthUser(h.listGithubRepositoryTags))
		})
	})

	h.Mux.Route("/projects/invitations", func(r chi.Router) {
		r.Use(middleware.HandleInvitationToken)
		r.Get("/accept", h.acceptInvitation)
		r.Get("/reject", h.rejectInvitation)
	})

	h.Mux.Route("/organization/settings", func(r chi.Router) {
		r.Get("/", middleware.RequireAuthUser(h.getSettings))
		r.Patch("/", middleware.RequireAuthUser(h.updateSettings))
	})

	h.Mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.Mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		util.WriteResponseError(w, responseerrors.NewNotFoundError())
	})
	h.Mux.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		util.WriteResponseError(w, responseerrors.NewMethodNotAllowedError())
	})
}
