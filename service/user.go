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

func (s *UserService) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return s.repository.ReadByEmail(ctx, email)
}
