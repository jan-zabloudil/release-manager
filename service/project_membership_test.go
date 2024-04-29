package service

import (
	"context"
	"testing"

	"release-manager/pkg/apierrors"
	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/dberrors"
	repo "release-manager/repository/mock"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProjectService_CreateInvitation(t *testing.T) {
	testCases := []struct {
		name      string
		creation  model.CreateProjectInvitationInput
		mockSetup func(*svc.AuthService, *svc.EmailService, *svc.ProjectService, *repo.ProjectInvitationRepository)
		wantErr   bool
	}{
		{
			name:     "Unknown project",
			creation: model.CreateProjectInvitationInput{},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, apierrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Invalid project invitation - missing email",
			creation: model.CreateProjectInvitationInput{
				Email:       "",
				ProjectRole: "editor",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Invalid project invitation - invalid role",
			creation: model.CreateProjectInvitationInput{
				Email:       "",
				ProjectRole: "new",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invitation already exists",
			creation: model.CreateProjectInvitationInput{
				Email:       "test@test.tt",
				ProjectRole: "viewer",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				repo.On("ReadByEmailForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Success",
			creation: model.CreateProjectInvitationInput{
				Email:       "test@test.tt",
				ProjectRole: "viewer",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				repo.On("ReadByEmailForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
				email.On("SendProjectInvitation", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(repo.ProjectInvitationRepository)
			project := new(svc.ProjectService)
			auth := new(svc.AuthService)
			invitationSender := new(svc.EmailService)
			service := NewProjectMembershipService(auth, project, repository, invitationSender)

			tc.mockSetup(auth, invitationSender, project, repository)

			_, err := service.CreateInvitation(context.Background(), tc.creation, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			project.AssertExpectations(t)
			repository.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetInvitations(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthService, *svc.ProjectService, *repo.ProjectInvitationRepository)
		wantErr   bool
	}{
		{
			name: "Unknown project",
			mockSetup: func(auth *svc.AuthService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				repo.On("ReadAllForProject", mock.Anything, mock.Anything).Return(
					[]model.ProjectInvitation{
						{Email: "test@test.tt", ProjectRole: model.ProjectRoleEditor, Status: model.InvitationStatusPending, ProjectID: uuid.New()},
					},
					nil,
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(repo.ProjectInvitationRepository)
			project := new(svc.ProjectService)
			auth := new(svc.AuthService)
			invitationSender := new(svc.EmailService)
			service := NewProjectMembershipService(auth, project, repository, invitationSender)

			tc.mockSetup(auth, project, repository)

			_, err := service.ListInvitations(context.Background(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			project.AssertExpectations(t)
			repository.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteInvitation(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthService, *svc.ProjectService, *repo.ProjectInvitationRepository)
		wantErr   bool
	}{
		{
			name: "Unknown project",
			mockSetup: func(auth *svc.AuthService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Unknown invitation",
			mockSetup: func(auth *svc.AuthService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				repo.On("Read", mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthService, project *svc.ProjectService, repo *repo.ProjectInvitationRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				project.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				repo.On("Read", mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, nil)
				repo.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(repo.ProjectInvitationRepository)
			project := new(svc.ProjectService)
			auth := new(svc.AuthService)
			invitationSender := new(svc.EmailService)
			service := NewProjectMembershipService(auth, project, repository, invitationSender)

			tc.mockSetup(auth, project, repository)

			err := service.DeleteInvitation(context.Background(), uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			project.AssertExpectations(t)
			repository.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_AcceptInvitation(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*repo.ProjectInvitationRepository)
		wantErr   bool
	}{
		{
			name: "Unknown invitation",
			mockSetup: func(invitationRepo *repo.ProjectInvitationRepository) {
				invitationRepo.On("ReadByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(invitationRepo *repo.ProjectInvitationRepository) {
				invitationRepo.On("ReadByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					model.ProjectInvitation{
						Email: "test@test.tt", ProjectRole: model.ProjectRoleEditor, Status: model.InvitationStatusPending, ProjectID: uuid.New(),
					},
					nil,
				)
				invitationRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(repo.ProjectInvitationRepository)
			project := new(svc.ProjectService)
			auth := new(svc.AuthService)
			invitationSender := new(svc.EmailService)
			service := NewProjectMembershipService(auth, project, repository, invitationSender)

			tc.mockSetup(repository)

			tkn, err := cryptox.NewToken()
			if err != nil {
				t.Fatal(err)
			}

			err = service.AcceptInvitation(context.Background(), tkn)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repository.AssertExpectations(t)
		})
	}
}

func TestProjectService_RejectInvitation(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*repo.ProjectInvitationRepository)
		wantErr   bool
	}{
		{
			name: "Unknown invitation",
			mockSetup: func(invitationRepo *repo.ProjectInvitationRepository) {
				invitationRepo.On("ReadByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(invitationRepo *repo.ProjectInvitationRepository) {
				invitationRepo.On("ReadByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					model.ProjectInvitation{
						Email: "test@test.tt", ProjectRole: model.ProjectRoleEditor, Status: model.InvitationStatusPending, ProjectID: uuid.New(),
					},
					nil,
				)
				invitationRepo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(repo.ProjectInvitationRepository)
			project := new(svc.ProjectService)
			auth := new(svc.AuthService)
			invitationSender := new(svc.EmailService)
			service := NewProjectMembershipService(auth, project, repository, invitationSender)

			tc.mockSetup(repository)

			tkn, err := cryptox.NewToken()
			if err != nil {
				t.Fatal(err)
			}

			err = service.RejectInvitation(context.Background(), tkn)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repository.AssertExpectations(t)
		})
	}
}
