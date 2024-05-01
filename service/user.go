package service

import (
	"context"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/dberrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type UserService struct {
	authGuard  authGuard
	repository userRepository
}

func NewUserService(guard authGuard, repo userRepository) *UserService {
	return &UserService{
		authGuard:  guard,
		repository: repo,
	}
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID, authUserID uuid.UUID) (model.User, error) {
	if err := s.authGuard.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return model.User{}, err
	}

	u, err := s.repository.Read(ctx, id)
	if err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return model.User{}, apierrors.NewUserNotFoundError().Wrap(err)
		default:
			return model.User{}, err
		}
	}

	return u, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return err
	}

	_, err := s.Get(ctx, id, authUserID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, id)
}

func (s *UserService) ListAll(ctx context.Context, authUserID uuid.UUID) ([]model.User, error) {
	if err := s.authGuard.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return nil, err
	}

	u, err := s.repository.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}
