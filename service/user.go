package service

import (
	"context"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/dberrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type UserService struct {
	authSvc    model.AuthService
	repository model.UserRepository
}

func NewUserService(authSvc model.AuthService, repo model.UserRepository) *UserService {
	return &UserService{
		authSvc:    authSvc,
		repository: repo,
	}
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID, authUserID uuid.UUID) (model.User, error) {
	if err := s.authSvc.AuthorizeAdminRole(ctx, authUserID); err != nil {
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
	if err := s.authSvc.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return err
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		switch {
		case dberrors.IsNotFoundError(err):
			return apierrors.NewUserNotFoundError().Wrap(err)
		default:
			return err
		}
	}

	return nil
}

func (s *UserService) GetAll(ctx context.Context, authUserID uuid.UUID) ([]model.User, error) {
	if err := s.authSvc.AuthorizeAdminRole(ctx, authUserID); err != nil {
		return nil, err
	}

	u, err := s.repository.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}
