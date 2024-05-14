package service

import (
	"context"

	"release-manager/pkg/apierrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type SettingsService struct {
	authGuard  authGuard
	repository settingsRepository
}

func NewSettingsService(guard authGuard, r settingsRepository) *SettingsService {
	return &SettingsService{
		authGuard:  guard,
		repository: r,
	}
}

func (s *SettingsService) Get(ctx context.Context, authUserID uuid.UUID) (model.Settings, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Settings{}, err
	}

	return s.repository.Read(ctx)
}

func (s *SettingsService) Update(ctx context.Context, u model.UpdateSettingsInput, authUserID uuid.UUID) (model.Settings, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Settings{}, err
	}

	settings, err := s.repository.Update(ctx, func(s model.Settings) (model.Settings, error) {
		if err := s.Update(u); err != nil {
			return model.Settings{}, apierrors.NewSettingsUnprocessableError().Wrap(err).WithMessage(err.Error())
		}

		return s, nil
	})
	if err != nil {
		return model.Settings{}, err
	}

	return settings, nil
}

func (s *SettingsService) GetGithubToken(ctx context.Context) (string, error) {
	settings, err := s.repository.Read(ctx)
	if err != nil {
		return "", err
	}

	if !settings.Github.Enabled {
		return "", apierrors.NewGithubIntegrationNotEnabledError()
	}

	return settings.Github.Token, nil
}
