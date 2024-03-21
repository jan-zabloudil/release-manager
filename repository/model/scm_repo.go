package model

import (
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type SCMRepo struct {
	Data Data `json:"scm_repo"`
}

type Data struct {
	Platform        string `json:"platform"`
	RepoURL         string `json:"repo_url"`
	OwnerIdentifier string `json:"owner_identifier"`
	RepoIdentifier  string `json:"repo_identifier"`
}

func ToDBSCMRepo(platform, repoUrl, ownerId, repoId string) SCMRepo {
	return SCMRepo{
		Data{
			Platform:        platform,
			RepoURL:         repoUrl,
			OwnerIdentifier: ownerId,
			RepoIdentifier:  repoId,
		},
	}
}

func ToSvcSCMRepo(appID uuid.UUID, platform string, repoURL string) (svcmodel.SCMRepo, error) {
	r, err := svcmodel.NewSCMRepo(appID, platform, repoURL)
	if err != nil {
		return nil, err
	}

	return r, nil
}
