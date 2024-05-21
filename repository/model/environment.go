package model

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Environment struct {
	ID         uuid.UUID `db:"id"`
	ProjectID  uuid.UUID `db:"project_id"`
	Name       string    `db:"name"`
	ServiceURL string    `db:"service_url"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	// release is optional, one or no release can be deployed to an environment
	ReleaseID        sql.NullString `db:"release_id"`
	ReleaseTitle     sql.NullString `db:"release_title"`
	ReleaseNotes     sql.NullString `db:"release_notes"`
	ReleaseCreatedBy sql.NullString `db:"release_created_by"`
	ReleaseCreatedAt sql.NullTime   `db:"release_created_at"`
	ReleaseUpdatedAt sql.NullTime   `db:"release_updated_at"`
}

func ToSvcEnvironment(e Environment) (svcmodel.Environment, error) {
	u, err := url.Parse(e.ServiceURL)
	if err != nil {
		return svcmodel.Environment{}, fmt.Errorf("parsing service URL: %w", err)
	}

	env := svcmodel.Environment{
		ID:         e.ID,
		ProjectID:  e.ProjectID,
		Name:       e.Name,
		ServiceURL: *u,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}

	if e.ReleaseID.Valid {
		rlsID, err := uuid.Parse(e.ReleaseID.String)
		if err != nil {
			return svcmodel.Environment{}, fmt.Errorf("parsing release ID: %w", err)
		}
		rlsCreatedByID, err := uuid.Parse(e.ReleaseCreatedBy.String)
		if err != nil {
			return svcmodel.Environment{}, fmt.Errorf("parsing release created by ID: %w", err)
		}

		rls := svcmodel.Release{
			ID:           rlsID,
			ProjectID:    e.ProjectID,
			ReleaseTitle: e.ReleaseTitle.String,
			ReleaseNotes: e.ReleaseNotes.String,
			AuthorUserID: rlsCreatedByID,
			CreatedAt:    e.ReleaseCreatedAt.Time,
			UpdatedAt:    e.ReleaseUpdatedAt.Time,
		}

		env.DeployedRelease = &rls
	}

	return env, nil
}

func ToSvcEnvironments(envs []Environment) ([]svcmodel.Environment, error) {
	svcEnvs := make([]svcmodel.Environment, 0, len(envs))
	for _, e := range envs {
		svcEnv, err := ToSvcEnvironment(e)
		if err != nil {
			return nil, err
		}

		svcEnvs = append(svcEnvs, svcEnv)
	}

	return svcEnvs, nil
}
