package service

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
)

type TemplateService struct {
	repository model.TemplateRepository
}

func (s *TemplateService) Create(ctx context.Context, t model.Template) (model.Template, error) {
	return s.repository.Insert(ctx, t)
}

func (s *TemplateService) GetAll(ctx context.Context) ([]model.Template, error) {
	return s.repository.ReadAll(ctx)
}

func (s *TemplateService) Get(ctx context.Context, id uuid.UUID) (model.Template, error) {
	return s.repository.Read(ctx, id)
}

func (s *TemplateService) Update(ctx context.Context, t model.Template) (model.Template, error) {
	return s.repository.Update(ctx, t)
}

func (s *TemplateService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
