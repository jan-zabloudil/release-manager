package service

import (
	"context"
	"fmt"

	svcerrors "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ReleaseService struct {
	authGuard      authGuard
	projectGetter  projectGetter
	settingsGetter settingsGetter
	slackNotifier  slackNotifier
	githubManager  githubManager
	repo           releaseRepository
}

func NewReleaseService(
	authGuard authGuard,
	projectGetter projectGetter,
	settingsGetter settingsGetter,
	notifier slackNotifier,
	manager githubManager,
	repo releaseRepository,
) *ReleaseService {
	return &ReleaseService{
		authGuard:      authGuard,
		projectGetter:  projectGetter,
		settingsGetter: settingsGetter,
		slackNotifier:  notifier,
		githubManager:  manager,
		repo:           repo,
	}
}

func (s *ReleaseService) Create(
	ctx context.Context,
	input model.CreateReleaseInput,
	projectID,
	authorUserID uuid.UUID,
) (model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authorUserID); err != nil {
		return model.Release{}, fmt.Errorf("authorizing project member: %w", err)
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return model.Release{}, fmt.Errorf("getting github token: %w", err)
	}

	p, err := s.projectGetter.GetProject(ctx, projectID, authorUserID)
	if err != nil {
		return model.Release{}, fmt.Errorf("getting project: %w", err)
	}

	if !p.IsGithubRepoSet() {
		return model.Release{}, svcerrors.NewGithubRepoNotSetForProjectError()
	}

	// Before checking if the tag exists, validate if git tag name was provided in order to avoid unnecessary API calls.
	if err := input.ValidateGitTagName(); err != nil {
		return model.Release{}, svcerrors.NewReleaseUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	tagExists, err := s.githubManager.TagExists(ctx, tkn, *p.GithubRepo, input.GitTagName)
	if err != nil {
		return model.Release{}, fmt.Errorf("checking if tag exists: %w", err)
	}
	if !tagExists {
		return model.Release{}, svcerrors.NewGitTagNotFoundError()
	}

	gitTagURL, err := s.githubManager.GenerateGitTagURL(*p.GithubRepo, input.GitTagName)
	if err != nil {
		return model.Release{}, fmt.Errorf("generating git tag URL: %w", err)
	}
	input.AddGitTagURL(gitTagURL)

	rls, err := model.NewRelease(input, projectID, authorUserID)
	if err != nil {
		return model.Release{}, svcerrors.NewReleaseUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.repo.Create(ctx, rls); err != nil {
		return model.Release{}, fmt.Errorf("creating release: %w", err)
	}

	return rls, nil
}

func (s *ReleaseService) Get(ctx context.Context, projectID, releaseID, authorUserID uuid.UUID) (model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authorUserID); err != nil {
		return model.Release{}, fmt.Errorf("authorizing project member: %w", err)
	}

	rls, err := s.repo.Read(ctx, projectID, releaseID)
	if err != nil {
		return model.Release{}, fmt.Errorf("reading release: %w", err)
	}

	return rls, nil
}

func (s *ReleaseService) Delete(ctx context.Context, projectID, releaseID, authorUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authorUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	if err := s.repo.Delete(ctx, projectID, releaseID); err != nil {
		return fmt.Errorf("deleting release: %w", err)
	}

	return nil
}

func (s *ReleaseService) Update(
	ctx context.Context,
	input model.UpdateReleaseInput,
	projectID,
	releaseID,
	authorUserID uuid.UUID,
) (model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authorUserID); err != nil {
		return model.Release{}, fmt.Errorf("authorizing project member: %w", err)
	}

	rls, err := s.repo.Update(ctx, projectID, releaseID, func(rls model.Release) (model.Release, error) {
		if err := rls.Update(input); err != nil {
			return model.Release{}, svcerrors.NewReleaseUnprocessableError().Wrap(err).WithMessage(err.Error())
		}

		return rls, nil
	})
	if err != nil {
		return model.Release{}, fmt.Errorf("updating release: %w", err)
	}

	return rls, nil
}

func (s *ReleaseService) ListForProject(ctx context.Context, projectID, authorUserID uuid.UUID) ([]model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authorUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	rls, err := s.repo.ListForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}
	if len(rls) == 0 {
		exists, err := s.projectGetter.ProjectExists(ctx, projectID, authorUserID)
		if err != nil {
			return nil, fmt.Errorf("checking project existence: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewProjectNotFoundError()
		}
	}

	return rls, nil
}

func (s *ReleaseService) SendReleaseNotification(ctx context.Context, projectID, releaseID, authorUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authorUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	rls, err := s.repo.Read(ctx, projectID, releaseID)
	if err != nil {
		return fmt.Errorf("reading release: %w", err)
	}

	p, err := s.projectGetter.GetProject(ctx, projectID, authorUserID)
	if err != nil {
		return fmt.Errorf("getting project: %w", err)
	}

	if !p.IsSlackChannelSet() {
		return svcerrors.NewSlackChannelNotSetForProjectError()
	}

	tkn, err := s.settingsGetter.GetSlackToken(ctx)
	if err != nil {
		return fmt.Errorf("getting slack token: %w", err)
	}

	if err := s.slackNotifier.SendReleaseNotification(ctx, tkn, p.SlackChannelID, model.NewReleaseNotification(p, rls)); err != nil {
		return fmt.Errorf("sending slack notification: %w", err)
	}

	return nil
}
