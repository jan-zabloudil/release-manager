package service

import (
	"context"
	"testing"

	repo "release-manager/repository/mock"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSettingsService_Update(t *testing.T) {
	settings := model.Settings{
		OrganizationName:      "Old Organization",
		DefaultReleaseMessage: "Old Message",
	}

	validName := "new name"
	invalidName := ""

	enabled := true
	validToken := "valid token"
	invalidToken := ""

	testCases := []struct {
		name      string
		userID    uuid.UUID
		update    model.UpdateSettingsInput
		mockSetup func(*svc.AuthService, *repo.SettingsRepository)
		expectErr bool
	}{
		{
			name:   "Success",
			userID: uuid.New(),
			update: model.UpdateSettingsInput{
				OrganizationName: &validName,
				Slack: model.UpdateSlackSettingsInput{
					Enabled: &enabled,
					Token:   &validToken,
				},
			},
			mockSetup: func(authSvc *svc.AuthService, settingsRepo *repo.SettingsRepository) {
				authSvc.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				settingsRepo.On("Read", mock.Anything, mock.Anything).Return(settings, nil)
				settingsRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectErr: false,
		},
		{
			name:   "Error - missing slack token",
			userID: uuid.New(),
			update: model.UpdateSettingsInput{
				OrganizationName: &validName,
				Slack: model.UpdateSlackSettingsInput{
					Enabled: &enabled,
					Token:   &invalidToken,
				},
			},
			mockSetup: func(authSvc *svc.AuthService, settingsRepo *repo.SettingsRepository) {
				authSvc.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				settingsRepo.On("Read", mock.Anything, mock.Anything).Return(settings, nil)
			},
			expectErr: true,
		},
		{
			name:   "Error - missing name",
			userID: uuid.New(),
			update: model.UpdateSettingsInput{
				OrganizationName: &invalidName,
			},
			mockSetup: func(authSvc *svc.AuthService, settingsRepo *repo.SettingsRepository) {
				authSvc.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				settingsRepo.On("Read", mock.Anything, mock.Anything).Return(settings, nil)
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthService)
			settingsRepo := new(repo.SettingsRepository)
			settingsSvc := NewSettingsService(authSvc, settingsRepo)

			tc.mockSetup(authSvc, settingsRepo)

			_, err := settingsSvc.Update(context.Background(), tc.update, tc.userID)

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
