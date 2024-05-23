package service

import (
	"context"
	"fmt"

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
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.User{}, fmt.Errorf("authorizing user role: %w", err)
	}

	u, err := s.repository.Read(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("reading user: %w", err)
	}

	return u, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := s.repository.ReadByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("reading user by email: %w", err)
	}

	return u, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return err
	}

	_, err := s.Get(ctx, id, authUserID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	return s.repository.Delete(ctx, id)
}

func (s *UserService) ListAll(ctx context.Context, authUserID uuid.UUID) ([]model.User, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}

	u, err := s.repository.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing all users: %w", err)
	}

	return u, nil
}
