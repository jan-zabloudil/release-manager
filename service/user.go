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

// Get retrieves a user by ID, can be accessed by the user themselves.
func (s *UserService) Get(ctx context.Context, userID uuid.UUID) (model.User, error) {
	if err := s.authGuard.AuthorizeUserRoleUser(ctx, userID); err != nil {
		return model.User{}, fmt.Errorf("authorizing user role: %w", err)
	}

	u, err := s.repository.Read(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("reading user: %w", err)
	}

	return u, nil
}

// GetByEmail retrieves a user by email, this function does not require authentication!
func (s *UserService) GetByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := s.repository.ReadByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("reading user by email: %w", err)
	}

	return u, nil
}

// GetForAdmin retrieves a user by ID, can be accessed only by admin user.
func (s *UserService) GetForAdmin(ctx context.Context, userID uuid.UUID, authUserID uuid.UUID) (model.User, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return model.User{}, fmt.Errorf("authorizing user role: %w", err)
	}

	u, err := s.repository.Read(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("reading user: %w", err)
	}

	return u, nil
}

// DeleteForAdmin deletes a user by ID, can be accessed only by admin user.
func (s *UserService) DeleteForAdmin(ctx context.Context, userID uuid.UUID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return err
	}

	_, err := s.GetForAdmin(ctx, userID, authUserID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	return s.repository.Delete(ctx, userID)
}

// ListAllForAdmin lists all users, can be accessed only by admin user.
func (s *UserService) ListAllForAdmin(ctx context.Context, authUserID uuid.UUID) ([]model.User, error) {
	if err := s.authGuard.AuthorizeUserRoleAdmin(ctx, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing user role: %w", err)
	}

	u, err := s.repository.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing all users: %w", err)
	}

	return u, nil
}
