package model

import (
	"net/url"
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Deployment struct {
	ID               id.Deployment `db:"id"`
	DeployedByUserID id.AuthUser   `db:"deployed_by"`
	DeployedAt       time.Time     `db:"deployed_at"`

	ReleaseID           uuid.UUID   `db:"release_id"`
	ReleaseProjectID    uuid.UUID   `db:"release_project_id"`
	ReleaseTitle        string      `db:"release_title"`
	ReleaseNotes        string      `db:"release_notes"`
	ReleaseAuthorUserID id.AuthUser `db:"release_created_by"`
	ReleaseCreatedAt    time.Time   `db:"release_created_at"`
	ReleaseUpdatedAt    time.Time   `db:"release_updated_at"`

	EnvID         id.Environment `db:"env_id"`
	EnvProjectID  uuid.UUID      `db:"env_project_id"`
	EnvName       string         `db:"env_name"`
	EnvServiceURL string         `db:"env_service_url"`
	EnvCreatedAt  time.Time      `db:"env_created_at"`
	EnvUpdatedAt  time.Time      `db:"env_updated_at"`
}

func ToSvcDeployment(dpl Deployment) (svcmodel.Deployment, error) {
	envURL, err := url.Parse(dpl.EnvServiceURL)
	if err != nil {
		return svcmodel.Deployment{}, err
	}

	return svcmodel.Deployment{
		ID:               dpl.ID,
		DeployedByUserID: dpl.DeployedByUserID,
		DeployedAt:       dpl.DeployedAt,
		Release: svcmodel.Release{
			ID:           dpl.ReleaseID,
			ProjectID:    dpl.ReleaseProjectID,
			ReleaseTitle: dpl.ReleaseTitle,
			ReleaseNotes: dpl.ReleaseNotes,
			AuthorUserID: dpl.ReleaseAuthorUserID,
			CreatedAt:    dpl.ReleaseCreatedAt,
			UpdatedAt:    dpl.ReleaseUpdatedAt,
		},
		Environment: svcmodel.Environment{
			ID:         dpl.EnvID,
			ProjectID:  dpl.EnvProjectID,
			Name:       dpl.EnvName,
			ServiceURL: *envURL,
			CreatedAt:  dpl.EnvCreatedAt,
			UpdatedAt:  dpl.EnvUpdatedAt,
		},
	}, nil
}

func ToSvcDeployments(dpls []Deployment) ([]svcmodel.Deployment, error) {
	d := make([]svcmodel.Deployment, 0, len(dpls))
	for _, dpl := range dpls {
		svcDpl, err := ToSvcDeployment(dpl)
		if err != nil {
			return nil, err
		}

		d = append(d, svcDpl)
	}

	return d, nil
}
