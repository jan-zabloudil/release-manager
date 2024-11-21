package service

import (
	"context"
	"testing"

	"release-manager/pkg/id"
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
	slackTkn := model.SlackToken("slackToken")
	emptySlackTkn := model.SlackToken("")

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
					Token:   &slackTkn,
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
					Token:   &emptySlackTkn,
				},
			},
			mockSetup: func(authSvc *svc.AuthorizationService, settingsRepo *repo.SettingsRepository) {
				authSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				settingsRepo.On("Upsert", mock.Anything, mock.Anything).Return(svcerrors.NewSettingsInvalidError())
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

			err := settingsSvc.Update(context.Background(), tc.update, id.AuthUser{})

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
