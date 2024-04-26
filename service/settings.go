package service

import (
	"context"

	"release-manager/pkg/apierrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type SettingsService struct {
	authService model.AuthService
	repository  model.SettingsRepository
}

func NewSettingsService(authSvc model.AuthService, r model.SettingsRepository) *SettingsService {
	return &SettingsService{
		authService: authSvc,
		repository:  r,
	}
}

func (s *SettingsService) Get(ctx context.Context, authUserID uuid.UUID) (model.Settings, error) {
	if err := s.authService.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return model.Settings{}, err
	}

	return s.repository.Read(ctx)
}

func (s *SettingsService) Update(ctx context.Context, u model.UpdateSettingsInput, authUserID uuid.UUID) (model.Settings, error) {
	if err := s.authService.AuthorizeAdminRole(ctx, authUserID); err != nil {
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
