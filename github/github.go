package github

import (
	"context"
	"log/slog"

	"release-manager/github/model"
	"release-manager/github/utils"
	svcmodel "release-manager/service/model"

	"github.com/google/go-github/v60/github"
)

const (
	tagsLimit = 20
)

// SilentGitHub is used when no GitHub API credentials are provided.
// SilentGithub does not call GitHub api
type SilentGitHub struct{}

type GitHub struct {
	client *github.Client
}

func New(token string) model.GitHub {
	if token == "" {
		return &SilentGitHub{}
	}

	return &GitHub{
		client: github.NewClient(nil).WithAuthToken(token),
	}
}

func (g *GitHub) ListTags(ctx context.Context, owner, repo string) ([]svcmodel.GitTag, error) {
	t, _, err := g.client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{PerPage: tagsLimit})
	if err != nil {
		return nil, utils.WrapGithubErr(err)
	}

	return model.ToSvcGitTags(t), nil
}

func (g *SilentGitHub) ListTags(ctx context.Context, owner, repo string) ([]svcmodel.GitTag, error) {
	slog.Debug("github api not called because silent github is used", "repo", repo)
	return nil, nil
}
