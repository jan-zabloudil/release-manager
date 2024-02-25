package service

import (
	"context"

	"release-manager/service/model"
)

type SettingsService struct {
	repository model.SettingsRepository
}

func (s *SettingsService) Set(ctx context.Context, sts model.Settings) (model.Settings, error) {
	return s.repository.Set(ctx, sts)
}

func (s *SettingsService) Get(ctx context.Context) (model.Settings, error) {
	return s.repository.Read(ctx)
}
