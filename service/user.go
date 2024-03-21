package service

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
)

type UserService struct {
	repository model.UserRepository
}

func (s *UserService) GetForToken(ctx context.Context, token string) (model.User, error) {
	return s.repository.ReadForToken(ctx, token)
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID) (model.User, error) {
	return s.repository.Read(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}

func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.repository.ReadAll(ctx)
}
