package service

import (
	"context"
	"fmt"

	svcerrors "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type ReleaseService struct {
	authGuard         authGuard
	projectGetter     projectGetter
	settingsGetter    settingsGetter
	environmentGetter environmentGetter
	slackNotifier     slackNotifier
	githubManager     githubManager
	repo              releaseRepository
}

func NewReleaseService(
	authGuard authGuard,
	projectGetter projectGetter,
	settingsGetter settingsGetter,
	environmentGetter environmentGetter,
	notifier slackNotifier,
	manager githubManager,
	repo releaseRepository,
) *ReleaseService {
	return &ReleaseService{
		authGuard:         authGuard,
		projectGetter:     projectGetter,
		settingsGetter:    settingsGetter,
		environmentGetter: environmentGetter,
		slackNotifier:     notifier,
		githubManager:     manager,
		repo:              repo,
	}
}

func (s *ReleaseService) CreateRelease(
	ctx context.Context,
	input model.CreateReleaseInput,
	projectID,
	authUserID uuid.UUID,
) (model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.Release{}, fmt.Errorf("authorizing project member: %w", err)
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return model.Release{}, fmt.Errorf("getting github token: %w", err)
	}

	p, err := s.projectGetter.GetProject(ctx, projectID, authUserID)
	if err != nil {
		return model.Release{}, fmt.Errorf("getting project: %w", err)
	}

	if !p.IsGithubRepoSet() {
		return model.Release{}, svcerrors.NewGithubRepoNotSetForProjectError()
	}

	if err := input.Validate(); err != nil {
		return model.Release{}, svcerrors.NewReleaseUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	tagExists, err := s.githubManager.TagExists(ctx, tkn, *p.GithubRepo, input.GitTagName)
	if err != nil {
		return model.Release{}, fmt.Errorf("checking if tag exists: %w", err)
	}
	if !tagExists {
		return model.Release{}, svcerrors.NewGitTagNotFoundError()
	}

	gitTagURL, err := s.githubManager.GenerateGitTagURL(p.GithubRepo.OwnerSlug, p.GithubRepo.RepoSlug, input.GitTagName)
	if err != nil {
		return model.Release{}, fmt.Errorf("generating git tag URL: %w", err)
	}

	rls, err := model.NewRelease(input, gitTagURL, projectID, authUserID)
	if err != nil {
		return model.Release{}, svcerrors.NewReleaseUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	if err := s.repo.CreateRelease(ctx, rls); err != nil {
		return model.Release{}, fmt.Errorf("creating release: %w", err)
	}

	return rls, nil
}

func (s *ReleaseService) GetRelease(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) (model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return model.Release{}, fmt.Errorf("authorizing project member: %w", err)
	}

	rls, err := s.repo.ReadRelease(ctx, projectID, releaseID)
	if err != nil {
		return model.Release{}, fmt.Errorf("reading release: %w", err)
	}

	return rls, nil
}

// DeleteRelease deletes a release. If deleteGithubRelease is true, it will also delete associacted GitHub release (if exists).
// Deleting GitHub release is idempotent, so if the release does not exist on GitHub, it will not return an error.
func (s *ReleaseService) DeleteRelease(ctx context.Context, input model.DeleteReleaseInput, projectID, releaseID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	if input.DeleteGithubRelease {
		err := s.deleteGithubRelease(ctx, projectID, releaseID, authUserID)
		if err != nil && !svcerrors.IsGithubReleaseNotFoundError(err) {
			return fmt.Errorf("deleting github release: %w", err)
		}
	}

	if err := s.repo.DeleteRelease(ctx, projectID, releaseID); err != nil {
		return fmt.Errorf("deleting release: %w", err)
	}

	return nil
}

// DeleteReleaseByGitTag is used when the release is deleted on GitHub and webhook is triggered to delete the release in the database.
func (s *ReleaseService) DeleteReleaseByGitTag(ctx context.Context, githubOwnerSlug, githubRepoSlug, gitTag string) error {
	if err := s.repo.DeleteReleaseByGitTag(ctx, githubOwnerSlug, githubRepoSlug, gitTag); err != nil {
		return fmt.Errorf("deleting release by git tag: %w", err)
	}

	return nil
}

func (s *ReleaseService) UpdateRelease(
	ctx context.Context,
	input model.UpdateReleaseInput,
	projectID,
	releaseID,
	authUserID uuid.UUID,
) (model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.Release{}, fmt.Errorf("authorizing project member: %w", err)
	}

	rls, err := s.repo.UpdateRelease(ctx, projectID, releaseID, func(rls model.Release) (model.Release, error) {
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

func (s *ReleaseService) ListReleasesForProject(ctx context.Context, projectID, authUserID uuid.UUID) ([]model.Release, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	rls, err := s.repo.ListReleasesForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}
	if len(rls) == 0 {
		exists, err := s.projectGetter.ProjectExists(ctx, projectID, authUserID)
		if err != nil {
			return nil, fmt.Errorf("checking project existence: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewProjectNotFoundError()
		}
	}

	return rls, nil
}

func (s *ReleaseService) SendReleaseNotification(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	tkn, err := s.settingsGetter.GetSlackToken(ctx)
	if err != nil {
		return fmt.Errorf("getting slack token: %w", err)
	}

	p, err := s.projectGetter.GetProject(ctx, projectID, authUserID)
	if err != nil {
		return fmt.Errorf("getting project: %w", err)
	}

	rls, err := s.repo.ReadRelease(ctx, projectID, releaseID)
	if err != nil {
		return fmt.Errorf("reading release: %w", err)
	}

	dpl, err := s.getLastDeploymentForRelease(ctx, projectID, releaseID)
	if err != nil {
		return fmt.Errorf("getting last deployment for release: %w", err)
	}

	if !p.IsSlackChannelSet() {
		return svcerrors.NewSlackChannelNotSetForProjectError()
	}

	if err := s.slackNotifier.SendReleaseNotification(ctx, tkn, p.SlackChannelID, model.NewReleaseNotification(p, rls, dpl)); err != nil {
		return fmt.Errorf("sending slack notification: %w", err)
	}

	return nil
}

func (s *ReleaseService) UpsertGithubRelease(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) error {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return fmt.Errorf("authorizing project member: %w", err)
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return fmt.Errorf("getting github token: %w", err)
	}

	p, err := s.projectGetter.GetProject(ctx, projectID, authUserID)
	if err != nil {
		return fmt.Errorf("getting project: %w", err)
	}

	rls, err := s.repo.ReadRelease(ctx, projectID, releaseID)
	if err != nil {
		return fmt.Errorf("reading release: %w", err)
	}

	if !p.IsGithubRepoSet() {
		return svcerrors.NewGithubRepoNotSetForProjectError()
	}

	if err := s.githubManager.UpsertRelease(ctx, tkn, *p.GithubRepo, rls); err != nil {
		return fmt.Errorf("upserting github release: %w", err)
	}

	return nil
}

func (s *ReleaseService) deleteGithubRelease(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) error {
	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return fmt.Errorf("getting github token: %w", err)
	}

	p, err := s.projectGetter.GetProject(ctx, projectID, authUserID)
	if err != nil {
		return fmt.Errorf("getting project: %w", err)
	}

	rls, err := s.repo.ReadRelease(ctx, projectID, releaseID)
	if err != nil {
		return fmt.Errorf("reading release: %w", err)
	}

	if !p.IsGithubRepoSet() {
		return svcerrors.NewGithubRepoNotSetForProjectError()
	}

	if err := s.githubManager.DeleteReleaseByTag(ctx, tkn, *p.GithubRepo, rls.GitTagName); err != nil {
		return fmt.Errorf("deleting github release: %w", err)
	}

	return nil
}

func (s *ReleaseService) GenerateGithubReleaseNotes(
	ctx context.Context,
	input model.GithubGeneratedReleaseNotesInput,
	projectID,
	authUserID uuid.UUID,
) (model.GithubGeneratedReleaseNotes, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.GithubGeneratedReleaseNotes{}, fmt.Errorf("authorizing project member: %w", err)
	}

	tkn, err := s.settingsGetter.GetGithubToken(ctx)
	if err != nil {
		return model.GithubGeneratedReleaseNotes{}, fmt.Errorf("getting Github token: %w", err)
	}

	project, err := s.projectGetter.GetProject(ctx, projectID, authUserID)
	if err != nil {
		return model.GithubGeneratedReleaseNotes{}, fmt.Errorf("reading project: %w", err)
	}

	if !project.IsGithubRepoSet() {
		return model.GithubGeneratedReleaseNotes{}, svcerrors.NewGithubRepoNotSetForProjectError()
	}

	if err := input.Validate(); err != nil {
		return model.GithubGeneratedReleaseNotes{}, svcerrors.NewGithubGeneratedNotesInvalidInputError().WithMessage(err.Error())
	}

	notes, err := s.githubManager.GenerateReleaseNotes(ctx, tkn, *project.GithubRepo, input)
	if err != nil {
		return model.GithubGeneratedReleaseNotes{}, fmt.Errorf("generating release notes: %w", err)
	}

	return notes, nil
}

func (s *ReleaseService) CreateDeployment(
	ctx context.Context,
	input model.CreateDeploymentInput,
	projectID,
	authUserID uuid.UUID,
) (model.Deployment, error) {
	if err := s.authGuard.AuthorizeProjectRoleEditor(ctx, projectID, authUserID); err != nil {
		return model.Deployment{}, fmt.Errorf("authorizing project member: %w", err)
	}

	if err := input.Validate(); err != nil {
		return model.Deployment{}, svcerrors.NewDeploymentUnprocessableError().Wrap(err).WithMessage(err.Error())
	}

	exists, err := s.projectGetter.ProjectExists(ctx, projectID, authUserID)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("checking if project exists: %w", err)
	}
	if !exists {
		return model.Deployment{}, svcerrors.NewProjectNotFoundError()
	}

	rls, err := s.repo.ReadRelease(ctx, projectID, input.ReleaseID)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("getting release: %w", err)
	}

	env, err := s.environmentGetter.GetEnvironment(ctx, projectID, input.EnvironmentID, authUserID)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("getting environment: %w", err)
	}

	dpl := model.NewDeployment(rls, env, authUserID)

	if err := s.repo.CreateDeployment(ctx, dpl); err != nil {
		return model.Deployment{}, fmt.Errorf("creating deployment: %w", err)
	}

	return dpl, nil
}

func (s *ReleaseService) ListDeploymentsForProject(
	ctx context.Context,
	params model.DeploymentFilterParams,
	projectID,
	authUserID uuid.UUID,
) ([]model.Deployment, error) {
	if err := s.authGuard.AuthorizeProjectRoleViewer(ctx, projectID, authUserID); err != nil {
		return nil, fmt.Errorf("authorizing project member: %w", err)
	}

	if params.ReleaseID != nil {
		exists, err := s.releaseExists(ctx, projectID, *params.ReleaseID)
		if err != nil {
			return nil, fmt.Errorf("checking if release exists: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewReleaseNotFoundError()
		}
	}

	if params.EnvironmentID != nil {
		exists, err := s.environmentGetter.EnvironmentExists(ctx, projectID, *params.EnvironmentID, authUserID)
		if err != nil {
			return nil, fmt.Errorf("checking if environment exists: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewEnvironmentNotFoundError()
		}
	}

	dpls, err := s.repo.ListDeploymentsForProject(ctx, params, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing deployments: %w", err)
	}

	if len(dpls) == 0 {
		exists, err := s.projectGetter.ProjectExists(ctx, projectID, authUserID)
		if err != nil {
			return nil, fmt.Errorf("checking if project exists: %w", err)
		}
		if !exists {
			return nil, svcerrors.NewProjectNotFoundError()
		}
	}

	return dpls, nil
}

// getLastDeploymentForRelease returns pointer to the last deployment for the release,
// or nil if no deployment exists for the release.
func (s *ReleaseService) getLastDeploymentForRelease(ctx context.Context, projectID, releaseID uuid.UUID) (*model.Deployment, error) {
	dpl, err := s.repo.ReadLastDeploymentForRelease(ctx, projectID, releaseID)
	if err != nil {
		switch {
		case svcerrors.IsNotFoundError(err):
			return nil, nil
		default:
			return nil, fmt.Errorf("reading last deployment for release: %w", err)
		}
	}

	return &dpl, nil
}

func (s *ReleaseService) releaseExists(ctx context.Context, projectID, releaseID uuid.UUID) (bool, error) {
	_, err := s.repo.ReadRelease(ctx, projectID, releaseID)
	if err != nil {
		switch {
		case svcerrors.IsNotFoundError(err):
			return false, nil
		default:
			return false, fmt.Errorf("checking if release exists: %w", err)
		}
	}

	return true, nil
}
