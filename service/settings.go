package service

import (
	"context"
	"fmt"

	svcerrors "release-manager/service/errors"
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
		return model.Settings{}, fmt.Errorf("authorizing user role: %w", err)
	}

	settings, err := s.repository.Read(ctx)
	if err != nil {
		return model.Settings{}, fmt.Errorf("reading settings: %w", err)
	}

	return settings, nil
}

func (s *SettingsService) Update(ctx context.Context, u model.UpdateSettingsInput, authUserID uuid.UUID) (model.Settings, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.Settings{}, fmt.Errorf("authorizing user role: %w", err)
	}

	settings, err := s.repository.Update(ctx, func(s model.Settings) (model.Settings, error) {
		if err := s.Update(u); err != nil {
			return model.Settings{}, svcerrors.NewSettingsUnprocessableError().Wrap(err).WithMessage(err.Error())
		}

		return s, nil
	})
	if err != nil {
		return model.Settings{}, fmt.Errorf("updating settings: %w", err)
	}

	return settings, nil
}

func (s *SettingsService) GetGithubToken(ctx context.Context) (string, error) {
	settings, err := s.repository.Read(ctx)
	if err != nil {
		return "", fmt.Errorf("reading settings: %w", err)
	}

	if !settings.Github.Enabled {
		return "", svcerrors.NewGithubIntegrationNotEnabledError()
	}

	return settings.Github.Token, nil
}

func (s *SettingsService) GetSlackToken(ctx context.Context) (string, error) {
	settings, err := s.repository.Read(ctx)
	if err != nil {
		return "", fmt.Errorf("reading settings: %w", err)
	}

	if !settings.Slack.Enabled {
		return "", svcerrors.NewSlackIntegrationNotEnabledError()
	}

	return settings.Slack.Token, nil
}
