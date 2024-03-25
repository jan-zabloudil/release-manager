package service

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
)

type ReleaseService struct {
	appSvc     model.AppService
	projectSvc model.ProjectService
	repository model.ReleaseRepository
	slack      model.Slack
}

func (s *ReleaseService) Create(ctx context.Context, r model.Release) (model.Release, error) {

	p, _ := s.projectSvc.Get(ctx, r.AppID)
	a, _ := s.appSvc.Get(ctx, r.AppID)

	s.slack.PostReleaseMessage(ctx, p, a, r)

	return s.repository.Insert(ctx, r)
}

func (s *ReleaseService) GetAllForApp(ctx context.Context, appID uuid.UUID) ([]model.Release, error) {
	return s.repository.ReadAllForApp(ctx, appID)
}

func (s *ReleaseService) Get(ctx context.Context, id uuid.UUID) (model.Release, error) {
	return s.repository.Read(ctx, id)
}

func (s *ReleaseService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}

func (s *ReleaseService) Update(ctx context.Context, rls model.Release) (model.Release, error) {
	return s.repository.Update(ctx, rls)
}
