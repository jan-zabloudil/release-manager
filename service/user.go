package service

import (
	"context"

	"release-manager/service/model"
)

type UserService struct {
	repository model.UserRepository
}

func (s *UserService) GetForToken(ctx context.Context, token string) (model.User, error) {
	return s.repository.ReadForToken(ctx, token)
}
