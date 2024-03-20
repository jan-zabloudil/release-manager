package service

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
)

type AppService struct {
	repository model.AppRepository
}

func (s *AppService) Create(ctx context.Context, app model.App) (model.App, error) {
	return s.repository.Insert(ctx, app)
}

func (s *AppService) Get(ctx context.Context, id uuid.UUID) (model.App, error) {
	return s.repository.Read(ctx, id)
}

func (s *AppService) GetAllForProject(ctx context.Context, projectID uuid.UUID) ([]model.App, error) {
	return s.repository.ReadAllForProject(ctx, projectID)
}

func (s *AppService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}

func (s *AppService) Update(ctx context.Context, app model.App) (model.App, error) {
	return s.repository.Update(ctx, app)
}
