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

	settings, err := s.Get(ctx, authUserID)
	if err != nil {
		return model.Settings{}, err
	}

	if err := settings.Update(u); err != nil {
		return model.Settings{}, apierrors.NewSettingsUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.repository.Update(ctx, settings); err != nil {
		return model.Settings{}, err
	}

	return settings, nil
}

func (s *SettingsService) GetGithubSettings(ctx context.Context) (model.GithubSettings, error) {
	settings, err := s.repository.Read(ctx)
	if err != nil {
		return model.GithubSettings{}, err
	}

	return settings.Github, nil
}
