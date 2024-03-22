package handler

import (
	"net/http"

	svcmodel "release-manager/service/model"
	neterrs "release-manager/transport/errors"
	"release-manager/transport/middleware"
	"release-manager/transport/utils"

	"github.com/go-chi/chi/v5"
	httpx "go.strv.io/net/http"
)

func (h *Handler) setupRoutes() {

	h.Mux.Use(httpx.RequestIDMiddleware(utils.RequestID))
	h.Mux.Use(httpx.LoggingMiddleware(utils.NewServerLogger("logging")))
	h.Mux.Use(httpx.RecoverMiddleware(utils.NewServerLogger("recover")))
	h.Mux.Use(middleware.Auth(h.UserSvc))

	h.Mux.Route("/projects", func(r chi.Router) {
		r.Post("/", middleware.RequireAdminUser(h.createProject))
		r.Get("/", middleware.RequireAuthUser(h.listProjects))
		r.Route("/{id}", func(r chi.Router) {
			r.Use(middleware.HandleProject(h.ProjectSvc))

			r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleViewer(), h.getProject))
			r.Patch("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleEditor(), h.updateProject))
			r.Delete("/", middleware.RequireAdminUser(h.deleteProject))

			r.Post("/memberships", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleEditor(), h.createProjectMembershipRequest)) // TODO better route url

			r.Route("/invitations", func(r chi.Router) {
				r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleViewer(), h.listProjectInvitations))
				r.Delete("/{invitationId}", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleEditor(), h.deleteProjectInvitation))
			})

			r.Route("/members", func(r chi.Router) {
				r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleViewer(), h.listProjectMembers))
				r.Route("/{userId}", func(r chi.Router) {
					r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleViewer(), h.handleProjectMember))
					r.Patch("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleEditor(), h.handleProjectMember))
					r.Delete("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleEditor(), h.handleProjectMember))
				})
			})

			r.Route("/apps", func(r chi.Router) {
				r.Post("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleEditor(), h.createApp))
				r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectID, svcmodel.ProjectRoleViewer(), h.listApps))
			})
		})
	})

	h.Mux.Route("/apps/{appId}", func(r chi.Router) {
		r.Use(middleware.HandleApp(h.AppSvc))
		r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleViewer(), h.getApp))
		r.Patch("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleEditor(), h.updateApp))
		r.Delete("/", middleware.RequireAdminUser(h.deleteApp))

		r.Route("/repository", func(r chi.Router) {
			r.Put("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleEditor(), h.setSCMRepo))
			r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleViewer(), h.getSCMRepo))
			r.Delete("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleEditor(), h.deleteSCMRepo))
			r.Get("/tags", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleViewer(), h.getSCMRepoTags))
		})

		r.Route("/releases", func(r chi.Router) {
			r.Post("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleEditor(), h.createRelease))
			r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleViewer(), h.listReleases))
		})
	})

	h.Mux.Route("/releases/{id}", func(r chi.Router) {
		r.Use(middleware.HandleRelease(h.ReleaseSvc, h.AppSvc))
		r.Get("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleViewer(), h.getRelease))
		r.Patch("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleEditor(), h.updateRelease))
		r.Delete("/", middleware.RequireProjectMemberRole(h.ProjectMemberSvc, utils.ContextProjectIDFromApp, svcmodel.ProjectRoleEditor(), h.deleteRelease))
	})

	h.Mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.Mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		utils.WriteNotFoundResponse(w, neterrs.ErrHttpNotFound)
	})
	h.Mux.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		utils.WriteMethodNotAllowedResponse(w, neterrs.ErrHttpMethodNotAllowed)
	})

}
