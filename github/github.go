package github

import (
	"context"
	"log/slog"

	"release-manager/github/model"
	"release-manager/github/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

type SilentGitHub struct{}

type GitHub struct {
	client   *github.Client
	ownerURL string
}

func New(authToken, owner string) *GitHub {
	return &GitHub{
		client:   github.NewClient(nil).WithAuthToken(authToken),
		ownerURL: owner,
	}
}

// NewSilent SilentGithub should be injected into services if no github token and owner url are provided.
// SilentGithub does not call github api
func NewSilent() *SilentGitHub {
	return &SilentGitHub{}
}

func (g *GitHub) CreateRelease(ctx context.Context, repo string, r svcmodel.GithubRelease) (svcmodel.GithubRelease, error) {
	rr := model.ToRepositoryRelease(
		r.TagName,
		r.TargetCommitish,
		r.Name,
		r.Changelog,
	)

	createdR, _, err := g.client.Repositories.CreateRelease(ctx, g.ownerURL, repo, &rr)
	if err != nil {
		return svcmodel.GithubRelease{}, utils.WrapGithubErr(err)
	}

	return model.ToSvcGithubRelease(
		createdR.TagName,
		createdR.TargetCommitish,
		createdR.Name,
		createdR.Body,
	), nil
}

func (g *GitHub) ListTags(ctx context.Context, repo string) ([]svcmodel.GitTag, error) {
	t, _, err := g.client.Repositories.ListTags(ctx, g.ownerURL, repo, &github.ListOptions{PerPage: 20})
	if err != nil {
		return nil, utils.WrapGithubErr(err)
	}

	return model.ToSvcGitTags(t), nil
}

func (g *SilentGitHub) CreateRelease(ctx context.Context, repo string, r svcmodel.GithubRelease) (svcmodel.GithubRelease, error) {
	slog.Debug("github api not called because silent github is used", "repo", repo)
	return svcmodel.GithubRelease{}, nil
}

func (g *SilentGitHub) ListTags(ctx context.Context, repo string) ([]svcmodel.GitTag, error) {
	slog.Debug("github api not called because silent github is used", "repo", repo)
	return nil, nil
}
