package handler

import (
	"context"

	cryptox "release-manager/pkg/crypto"
	svcmodel "release-manager/service/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, c svcmodel.CreateProjectInput, authUserID uuid.UUID) (svcmodel.Project, error)
	Get(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (svcmodel.Project, error)
	ListAll(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.Project, error)
	Update(ctx context.Context, u svcmodel.UpdateProjectInput, projectID, authUserID uuid.UUID) (svcmodel.Project, error)
	Delete(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) error

	CreateEnvironment(ctx context.Context, c svcmodel.CreateEnvironmentInput, authUserID uuid.UUID) (svcmodel.Environment, error)
	GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)
	ListEnvironments(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.Environment, error)
	DeleteEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) error
	UpdateEnvironment(ctx context.Context, u svcmodel.UpdateEnvironmentInput, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)
}

type ProjectMembershipService interface {
	CreateInvitation(ctx context.Context, c svcmodel.CreateProjectInvitationInput, authUserID uuid.UUID) (svcmodel.ProjectInvitation, error)
	ListInvitations(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.ProjectInvitation, error)
	DeleteInvitation(ctx context.Context, projectID, invitationID, authUserID uuid.UUID) error
	AcceptInvitation(ctx context.Context, tkn cryptox.Token) error
	RejectInvitation(ctx context.Context, tkn cryptox.Token) error
}

type UserService interface {
	Get(ctx context.Context, id, authUserID uuid.UUID) (svcmodel.User, error)
	ListAll(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.User, error)
	Delete(ctx context.Context, id, authUserID uuid.UUID) error
}

type AuthService interface {
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}

type SettingsService interface {
	Update(ctx context.Context, u svcmodel.UpdateSettingsInput, authUserID uuid.UUID) (svcmodel.Settings, error)
	Get(ctx context.Context, authUserID uuid.UUID) (svcmodel.Settings, error)
}

type Handler struct {
	Mux                  *chi.Mux
	AuthSvc              AuthService
	UserSvc              UserService
	ProjectSvc           ProjectService
	SettingsSvc          SettingsService
	ProjectMembershipSvc ProjectMembershipService
}

func NewHandler(
	as AuthService,
	us UserService,
	ps ProjectService,
	ss SettingsService,
	pms ProjectMembershipService,
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
