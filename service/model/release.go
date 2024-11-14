package model

import (
	"errors"
	"net/url"
	"time"

	"release-manager/pkg/id"

	"github.com/google/uuid"
)

var (
	errReleaseTitleRequired               = errors.New("release title is required")
	errGithubGeneratedNotesGitTagRequired = errors.New("git tag name is required")
)

type CreateReleaseInput struct {
	ReleaseTitle string
	ReleaseNotes string
	// Used for linking the release with a specific point in a git repository.
	GitTagName string
}

type UpdateReleaseInput struct {
	ReleaseTitle *string
	ReleaseNotes *string
}

type Release struct {
	ID           id.Release
	ProjectID    id.Project
	ReleaseTitle string
	ReleaseNotes string
	AuthorUserID id.AuthUser
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tag          GitTag
	Attachments  []ReleaseAttachment
}

type GitTag struct {
	Name string
	URL  url.URL
}

type ReleaseAttachment struct {
	ID        uuid.UUID
	Name      string
	FilePath  string
	URL       url.URL
	CreatedAt time.Time
}

func NewRelease(input CreateReleaseInput, tag GitTag, projectID id.Project, authorUserID id.AuthUser) (Release, error) {
	now := time.Now()
	r := Release{
		ID:           id.NewRelease(),
		ProjectID:    projectID,
		ReleaseTitle: input.ReleaseTitle,
		ReleaseNotes: input.ReleaseNotes,
		Tag:          tag,
		AuthorUserID: authorUserID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := r.Validate(); err != nil {
		return Release{}, err
	}

	return r, nil
}

func (r *Release) Update(input UpdateReleaseInput) error {
	if input.ReleaseTitle != nil {
		r.ReleaseTitle = *input.ReleaseTitle
	}
	if input.ReleaseNotes != nil {
		r.ReleaseNotes = *input.ReleaseNotes
	}

	r.UpdatedAt = time.Now()

	return r.Validate()
}

func (r *Release) Validate() error {
	if r.ReleaseTitle == "" {
		return errReleaseTitleRequired
	}

	return nil
}

type ReleaseNotification struct {
	Message               string
	ProjectName           *string
	ReleaseTitle          *string
	ReleaseNotes          *string
	GitTagName            *string
	GitTagURL             *url.URL
	DeployedToEnvironment *string
	DeployedAt            *time.Time
	DeployedServiceURL    *url.URL
}

func NewReleaseNotification(p Project, r Release, dpl *Deployment) ReleaseNotification {
	n := ReleaseNotification{
		Message: p.ReleaseNotificationConfig.Message,
	}

	if p.ReleaseNotificationConfig.ShowProjectName {
		n.ProjectName = &p.Name
	}
	if p.ReleaseNotificationConfig.ShowReleaseTitle {
		n.ReleaseTitle = &r.ReleaseTitle
	}
	if p.ReleaseNotificationConfig.ShowReleaseNotes {
		n.ReleaseNotes = &r.ReleaseNotes
	}
	if p.ReleaseNotificationConfig.ShowSourceCode {
		n.GitTagName = &r.Tag.Name
		n.GitTagURL = &r.Tag.URL
	}
	if p.ReleaseNotificationConfig.ShowLastDeployment && dpl != nil {
		n.DeployedToEnvironment = &dpl.Environment.Name
		n.DeployedAt = &dpl.DeployedAt

		// Add the service URL to notification only if it is set
		// Service URL is not required for all environments
		if dpl.Environment.IsServiceURLSet() {
			n.DeployedServiceURL = &dpl.Environment.ServiceURL
		}
	}

	return n
}

type DeleteReleaseInput struct {
	DeleteGithubRelease bool
}

type GithubReleaseNotesInput struct {
	GitTagName         *string
	PreviousGitTagName *string
}

func (i GithubReleaseNotesInput) Validate() error {
	if i.GitTagName == nil {
		return errGithubGeneratedNotesGitTagRequired
	}
	if *i.GitTagName == "" {
		return errGithubGeneratedNotesGitTagRequired
	}

	return nil
}

func (i GithubReleaseNotesInput) GetGitTagName() string {
	if i.GitTagName != nil {
		return *i.GitTagName
	}

	return ""
}

type GithubReleaseNotes struct {
	Title string
	Notes string
}

type GithubTagDeletionWebhookInput struct {
	RawPayload []byte
	Signature  string
}
