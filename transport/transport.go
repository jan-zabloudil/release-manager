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
	Mux                      *chi.Mux
	UserSvc                  model.UserService
	ProjectSvc               model.ProjectService
	ProjectMembershipMgmtSvc model.ProjectMembershipManagementService
	ProjectInvitationSvc     model.ProjectInvitationService
	ProjectMemberSvc         model.ProjectMemberService
}

func NewHandler(us model.UserService, ps model.ProjectService, pms model.ProjectMembershipManagementService, pis model.ProjectInvitationService, pmrs model.ProjectMemberService) *Handler {
	h := Handler{
		Mux:                      chi.NewRouter(),
		UserSvc:                  us,
		ProjectSvc:               ps,
		ProjectMembershipMgmtSvc: pms,
		ProjectInvitationSvc:     pis,
		ProjectMemberSvc:         pmrs,
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

			r.Post("/memberships", h.createProjectMembershipRequest) // TODO better route url
			r.Get("/invitations", h.listProjectInvitations)
			r.Route("/members", func(r chi.Router) {
				r.Get("/", h.listProjectMembers)
				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", h.handleProjectMember)
					r.Patch("/", h.handleProjectMember)
					r.Delete("/", h.handleProjectMember)
				})
			})
		})
	})

	h.Mux.Route("/invitations/{id}", func(r chi.Router) {
		// TODO implement requireProjectMembership logic
		r.Use(h.requireProjectMembership)
		r.Delete("/", h.deleteProjectInvitation)
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
