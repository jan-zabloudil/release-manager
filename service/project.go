package service

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectService struct {
	repository model.ProjectRepository
}

func (s *ProjectService) Create(ctx context.Context, p model.Project, userID uuid.UUID) (model.Project, error) {
	return s.repository.Insert(ctx, p, userID)
}

func (s *ProjectService) Get(ctx context.Context, id uuid.UUID) (model.Project, error) {
	return s.repository.Read(ctx, id)
}

func (s *ProjectService) ListAll(ctx context.Context, u model.User) ([]model.Project, error) {
	if u.IsAdmin {
		return s.repository.ReadAll(ctx)
	}

	return s.repository.ReadForUser(ctx, u.ID)
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}

func (s *ProjectService) Update(ctx context.Context, p model.Project) (model.Project, error) {
	return s.repository.Update(ctx, p)
}
