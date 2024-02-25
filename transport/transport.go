package transport

import (
	"net/http"

	neterrs "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"

	"github.com/go-chi/chi/v5"
	httpx "go.strv.io/net/http"
)

type Handler struct {
	Mux        *chi.Mux
	UserSvc    model.UserService
	ProjectSvc model.ProjectService
}

func NewHandler(us model.UserService, ps model.ProjectService) *Handler {
	h := Handler{
		Mux:        chi.NewRouter(),
		UserSvc:    us,
		ProjectSvc: ps,
	}

	h.Mux.Use(httpx.RequestIDMiddleware(RequestID))
	h.Mux.Use(httpx.LoggingMiddleware(utils.NewServerLogger("logging")))
	h.Mux.Use(httpx.RecoverMiddleware(utils.NewServerLogger("recover")))
	h.Mux.Use(h.auth)

	h.Mux.Route("/projects", func(r chi.Router) {
		r.Post("/", h.requireAdminUser(h.createProject))
		r.Get("/", h.requireAuthUser(h.listProjects))
		r.Route("/{id}", func(r chi.Router) {
			// TODO implement requireProjectMembership logic
			r.Use(h.requireProjectMembership)
			r.Use(h.handleProject)
			r.Get("/", h.getProject)
			r.Patch("/", h.updateProject)
			r.Delete("/", h.deleteProject)
		})
	})

	h.Mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.Mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		WriteNotFoundResponse(w, neterrs.ErrHttpNotFound)
	})
	h.Mux.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		WriteMethodNotAllowedResponse(w, neterrs.ErrHttpMethodNotAllowed)
	})

	return &h
}
