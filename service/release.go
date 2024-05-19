package service

import (
	"context"
	"log/slog"

	"release-manager/pkg/apierrors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ReleaseService struct {
	projectGetter  projectGetter
	settingsGetter settingsGetter
	slackNotifier  slackNotifier
	repo           releaseRepository
}

func NewReleaseService(
	projectGetter projectGetter,
	settingsGetter settingsGetter,
	notifier slackNotifier,
	repo releaseRepository,
) *ReleaseService {
	return &ReleaseService{
		projectGetter:  projectGetter,
		settingsGetter: settingsGetter,
		slackNotifier:  notifier,
		repo:           repo,
	}
}

func (s *ReleaseService) Create(
	ctx context.Context,
	input model.CreateReleaseInput,
	sendReleaseNotification bool,
	projectID,
	authorUserID uuid.UUID,
) (model.Release, error) {
	// TODO add project member authorization

	p, err := s.projectGetter.GetProject(ctx, projectID, authorUserID)
	if err != nil {
		return model.Release{}, err
	}

	rls, err := model.NewRelease(input, projectID, authorUserID)
	if err != nil {
		return model.Release{}, apierrors.NewReleaseUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	unique, err := s.isReleaseTitleUniqueInProject(ctx, projectID, rls.ReleaseTitle)
	if err != nil {
		return model.Release{}, err
	}
	if !unique {
		return model.Release{}, apierrors.NewReleaseDuplicateTitleError().Wrap(err)
	}

	if err := s.repo.Create(ctx, rls); err != nil {
		return model.Release{}, err
	}

	if sendReleaseNotification {
		s.sendReleaseNotification(ctx, p, rls)
	}

	return rls, nil
}

func (s *ReleaseService) Get(ctx context.Context, projectID, releaseID, authorUserID uuid.UUID) (model.Release, error) {
	// TODO add project member authorization

	rls, err := s.repo.Read(ctx, projectID, releaseID)
	if err != nil {
		return model.Release{}, err
	}

	return rls, nil
}

func (s *ReleaseService) Delete(ctx context.Context, projectID, releaseID, authorUserID uuid.UUID) error {
	// TODO add project member authorization

	if err := s.repo.Delete(ctx, projectID, releaseID); err != nil {
		return err
	}

	return nil
}

func (s *ReleaseService) ListForProject(ctx context.Context, projectID, authorUserID uuid.UUID) ([]model.Release, error) {
	// TODO add project member authorization

	rls, err := s.repo.ListForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(rls) == 0 {
		exists, err := s.projectGetter.ProjectExists(ctx, projectID, authorUserID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apierrors.NewProjectNotFoundError()
		}
	}

	return rls, nil
}

func (s *ReleaseService) sendReleaseNotification(ctx context.Context, p model.Project, rls model.Release) {
	if !p.IsSlackChannelSet() {
		slog.Debug("notification not sent: slack channel missing for project", "project_id", p.ID)
		return
	}

	tkn, err := s.settingsGetter.GetSlackToken(ctx)
	if err != nil {
		// when fetching slack token fails, just log the error
		// fail attempt to send Slack notification should not affect the release creation
		// two possible reasons for failure:
		// 1. slack integration is not set (logged in debug level, as it's not an error, but possible scenario)
		// 2. failed to fetch slack token (logged in error level)
		slog.Log(ctx, apierrors.GetLogLevel(err), "failed to get slack token", "err", err)
		return
	}

	s.slackNotifier.SendReleaseNotificationAsync(ctx, tkn, p.SlackChannelID, model.NewReleaseNotification(p, rls))
}

func (s *ReleaseService) isReleaseTitleUniqueInProject(ctx context.Context, projectID uuid.UUID, title string) (bool, error) {
	_, err := s.repo.ReadByTitle(ctx, projectID, title)
	if err != nil {
		if apierrors.IsNotFoundError(err) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
