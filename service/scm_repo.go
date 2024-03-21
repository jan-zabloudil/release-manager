package service

import (
	"context"

	svcerr "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type SCMRepoService struct {
	repository model.SCMRepoRepository
	github     model.GitHub
}

func (s *SCMRepoService) SetRepo(ctx context.Context, repo model.SCMRepo) (model.SCMRepo, error) {
	return s.repository.InsertRepo(ctx, repo)
}

func (s *SCMRepoService) GetRepo(ctx context.Context, appID uuid.UUID) (model.SCMRepo, error) {
	return s.repository.ReadRepo(ctx, appID)
}

func (s *SCMRepoService) DeleteRepo(ctx context.Context, appID uuid.UUID) error {
	return s.repository.DeleteRepo(ctx, appID)
}

func (s *SCMRepoService) GetTags(ctx context.Context, appID uuid.UUID) ([]model.GitTag, error) {
	repo, err := s.GetRepo(ctx, appID)
	if err != nil {
		return nil, err
	}

	if !repo.IsSet() {
		return nil, svcerr.ErrSCMRepoNotSet
	}

	// Currently only GitHub repositories are supported
	// It's reasonable to assume that the repository type is GitHub, as no other values are permissible.
	t, err := s.github.ListTags(ctx, repo.RepoOwnerIdentifier(), repo.RepoIdentifier())
	if err != nil {
		return nil, err
	}

	return t, nil
}
