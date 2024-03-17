package handler

import (
	"release-manager/transport/model"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Mux                      *chi.Mux
	UserSvc                  model.UserService
	ProjectSvc               model.ProjectService
	ProjectMembershipMgmtSvc model.ProjectMembershipManagementService
	ProjectInvitationSvc     model.ProjectInvitationService
	ProjectMemberSvc         model.ProjectMemberService
}

func NewHandler(
	us model.UserService,
	ps model.ProjectService,
	pms model.ProjectMembershipManagementService,
	pis model.ProjectInvitationService,
	pmrs model.ProjectMemberService,
) *Handler {
	h := Handler{
		Mux:                      chi.NewRouter(),
		UserSvc:                  us,
		ProjectSvc:               ps,
		ProjectMembershipMgmtSvc: pms,
		ProjectInvitationSvc:     pis,
		ProjectMemberSvc:         pmrs,
	}

	h.setupRoutes()

	return &h
}
