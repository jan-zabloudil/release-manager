package transport

import (
	"net/http"

	svcmodel "release-manager/service/model"
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
			r.Use(h.handleProject)

			r.Get("/", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleViewer(), h.getProject))
			r.Patch("/", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleEditor(), h.updateProject))
			r.Delete("/", h.requireAdminUser(h.deleteProject))

			r.Post("/memberships", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleEditor(), h.createProjectMembershipRequest)) // TODO better route url

			r.Route("/invitations", func(r chi.Router) {
				r.Get("/", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleViewer(), h.listProjectInvitations))
				r.Delete("/{invitationId}", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleEditor(), h.deleteProjectInvitation))
			})

			r.Route("/members", func(r chi.Router) {
				r.Get("/", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleViewer(), h.listProjectMembers))
				r.Route("/{userId}", func(r chi.Router) {
					r.Get("/", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleViewer(), h.handleProjectMember))
					r.Patch("/", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleEditor(), h.handleProjectMember))
					r.Delete("/", h.requireProjectMemberRole(ContextProjectID, svcmodel.ProjectRoleEditor(), h.handleProjectMember))
				})
			})
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
