package model

import (
	"context"
	"net/url"
	"strings"

	svcerr "release-manager/service/errors"

	"github.com/google/uuid"
)

const (
	githubPlatform     = "github"
	githubPlatformHost = "github.com"
)

type SCMRepoRepository interface {
	InsertRepo(ctx context.Context, repo SCMRepo) (SCMRepo, error)
	ReadRepo(ctx context.Context, appID uuid.UUID) (SCMRepo, error)
	DeleteRepo(ctx context.Context, appID uuid.UUID) error
}

type GitHub interface {
	ListTags(ctx context.Context, owner, repo string) ([]GitTag, error)
}

type SCMRepo interface {
	AppID() uuid.UUID
	Platform() string
	RepoAbsURL() string
	RepoOwnerIdentifier() string
	RepoIdentifier() string
	IsSet() bool
}

type githubRepo struct {
	appID     uuid.UUID
	platform  string
	repoURL   string
	ownerSlug string
	repoSlug  string
}

type emptyRepo struct{}

func NewSCMRepo(appID uuid.UUID, platform, repoURL string) (SCMRepo, error) {
	switch platform {
	case "":
		return &emptyRepo{}, nil
	case githubPlatform:
		return newGithubRepo(appID, repoURL)
	default:
		return nil, svcerr.ErrUnknownSCMRepoPlatform
	}
}

func NewEmptySCMRepo() SCMRepo {
	return &emptyRepo{}
}

func newGithubRepo(appID uuid.UUID, repoURL string) (SCMRepo, error) {
	parsedURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return nil, svcerr.ErrInvalidSCMRepoURL
	}
	if !parsedURL.IsAbs() {
		return nil, svcerr.ErrInvalidSCMRepoURL
	}

	if !strings.Contains(parsedURL.Host, githubPlatformHost) {
		return nil, svcerr.ErrInvalidGithubHostUrl
	}

	path := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(path) != 2 {
		return nil, svcerr.ErrInvalidGithubRepoUrl
	}

	ownerSlug, repoSlug := path[0], path[1]

	return &githubRepo{
		appID:     appID,
		platform:  githubPlatform,
		repoURL:   repoURL,
		ownerSlug: ownerSlug,
		repoSlug:  repoSlug,
	}, nil
}

func (r *githubRepo) AppID() uuid.UUID {
	return r.appID
}
func (r *githubRepo) Platform() string {
	return r.platform
}
func (r *githubRepo) RepoAbsURL() string {
	return r.repoURL
}
func (r *githubRepo) RepoOwnerIdentifier() string {
	return r.ownerSlug
}
func (r *githubRepo) RepoIdentifier() string {
	return r.repoSlug
}
func (r *githubRepo) IsSet() bool {
	return true
}

func (r *emptyRepo) IsSet() bool {
	return false
}
func (r *emptyRepo) AppID() uuid.UUID {
	return uuid.Nil
}
func (r *emptyRepo) Platform() string {
	return ""
}
func (r *emptyRepo) RepoAbsURL() string {
	return ""
}
func (r *emptyRepo) RepoOwnerIdentifier() string {
	return ""
}
func (r *emptyRepo) RepoIdentifier() string {
	return ""
}
