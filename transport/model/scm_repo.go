package model

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type SCMRepoService interface {
	SetRepo(ctx context.Context, repo svcmodel.SCMRepo) (svcmodel.SCMRepo, error)
	GetRepo(ctx context.Context, appID uuid.UUID) (svcmodel.SCMRepo, error)
	DeleteRepo(ctx context.Context, appID uuid.UUID) error
	GetTags(ctx context.Context, appID uuid.UUID) ([]svcmodel.GitTag, error)
}

type SCMRepo struct {
	Platform string `json:"platform" validate:"required"`
	RepoURL  string `json:"repo_url" validate:"required"`
}

func ToSvcSCMRepo(appID uuid.UUID, platform, repoURL string) (svcmodel.SCMRepo, error) {
	repo, err := svcmodel.NewSCMRepo(appID, platform, repoURL)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func ToNetSCMRepo(platform, repoUrl string) SCMRepo {
	return SCMRepo{
		Platform: platform,
		RepoURL:  repoUrl,
	}
}
