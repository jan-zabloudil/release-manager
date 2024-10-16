package service

import (
	"context"
	"testing"

	"release-manager/pkg/pointer"
	repo "release-manager/repository/mock"
	svcerrors "release-manager/service/errors"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSettingsService_Update(t *testing.T) {
	testCases := []struct {
		name      string
		userID    uuid.UUID
		update    model.UpdateSettingsInput
		mockSetup func(*svc.AuthorizationService, *repo.SettingsRepository)
		expectErr bool
	}{
		{
			name:   "Success",
			userID: uuid.New(),
			update: model.UpdateSettingsInput{
				OrganizationName: pointer.StringPtr("New Organization"),
				Slack: model.UpdateSlackSettingsInput{
					Enabled: pointer.BoolPtr(true),
					Token:   pointer.StringPtr("newToken"),
				},
			},
			mockSetup: func(authSvc *svc.AuthorizationService, settingsRepo *repo.SettingsRepository) {
				authSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				settingsRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)
			},
			expectErr: false,
		},
		{
			name:   "Error - invalid update input",
			userID: uuid.New(),
			update: model.UpdateSettingsInput{
				OrganizationName: pointer.StringPtr("New Organization"),
				Slack: model.UpdateSlackSettingsInput{
					Enabled: pointer.BoolPtr(true),
					Token:   pointer.StringPtr(""),
				},
			},
			mockSetup: func(authSvc *svc.AuthorizationService, settingsRepo *repo.SettingsRepository) {
				authSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				settingsRepo.On("Upsert", mock.Anything, mock.Anything).Return(svcerrors.NewSettingsUnprocessableError())
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizationService)
			settingsRepo := new(repo.SettingsRepository)
			settingsSvc := NewSettingsService(authSvc, settingsRepo)

			tc.mockSetup(authSvc, settingsRepo)

			err := settingsSvc.Update(context.Background(), tc.update, tc.userID)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			settingsRepo.AssertExpectations(t)
		})
	}
}
