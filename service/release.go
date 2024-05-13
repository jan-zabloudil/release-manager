package service

import (
	"context"

	"release-manager/pkg/apierrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ReleaseService struct {
	projectGetter projectGetter
	repo          releaseRepository
}

func NewReleaseService(projectGetter projectGetter, repo releaseRepository) *ReleaseService {
	return &ReleaseService{
		projectGetter: projectGetter,
		repo:          repo,
	}
}

func (s *ReleaseService) Create(ctx context.Context, input model.CreateReleaseInput, projectID, authorUserID uuid.UUID) (model.Release, error) {
	// TODO add project member authorization

	// More features are going to be added, project object will be needed here, therefore GetProject is called here (instead of projectExists)
	_, err := s.projectGetter.GetProject(ctx, projectID, authorUserID)
	if err != nil {
		return model.Release{}, err
	}

	rls, err := model.NewRelease(input, projectID, authorUserID)
	if err != nil {
		return model.Release{}, apierrors.NewReleaseUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	// TODO check if release name is unique per project
	// TODO send slack notification
	if err := s.repo.Create(ctx, rls); err != nil {
		return model.Release{}, err
	}

	return rls, nil
}
