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
	AppSvc                   model.AppService
}

func NewHandler(
	us model.UserService,
	ps model.ProjectService,
	pms model.ProjectMembershipManagementService,
	pis model.ProjectInvitationService,
	pmrs model.ProjectMemberService,
	as model.AppService,
) *Handler {
	h := Handler{
		Mux:                      chi.NewRouter(),
		UserSvc:                  us,
		ProjectSvc:               ps,
		ProjectMembershipMgmtSvc: pms,
		ProjectInvitationSvc:     pis,
		ProjectMemberSvc:         pmrs,
		AppSvc:                   as,
	}

	h.setupRoutes()

	return &h
}
