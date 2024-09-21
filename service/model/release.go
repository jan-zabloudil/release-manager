package model

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

var (
	errReleaseTitleRequired = errors.New("release title is required")

	errReleaseGitTagRequired              = errors.New("git tag name is required")
	errReleaseGitTagURLRequired           = errors.New("git tag URL is required")
	errGithubGeneratedNotesGitTagRequired = errors.New("git tag name is required")

	errReleaseAttachmentNameRequired     = errors.New("attachment name is required")
	errReleaseAttachmentFilePathRequired = errors.New("attachment file path is required")
	errReleaseAttachmentURLRequired      = errors.New("attachment URL cannot be empty")
)

type CreateReleaseInput struct {
	ReleaseTitle string
	ReleaseNotes string
	// Used for linking the release with a specific point in a git repository.
	GitTagName string
}

func (i *CreateReleaseInput) Validate() error {
	if i.ReleaseTitle == "" {
		return errReleaseTitleRequired
	}
	if i.GitTagName == "" {
		return errReleaseGitTagRequired
	}

	return nil
}

type UpdateReleaseInput struct {
	ReleaseTitle *string
	ReleaseNotes *string
}

type Release struct {
	ID           uuid.UUID
	ProjectID    uuid.UUID
	ReleaseTitle string
	ReleaseNotes string
	AuthorUserID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	GitTagName   string
	GitTagURL    url.URL
	Attachments  []ReleaseAttachment
}

func NewRelease(input CreateReleaseInput, tagURL url.URL, projectID, authorUserID uuid.UUID) (Release, error) {
	now := time.Now()
	r := Release{
		ID:           uuid.New(),
		ProjectID:    projectID,
		ReleaseTitle: input.ReleaseTitle,
		ReleaseNotes: input.ReleaseNotes,
		GitTagName:   input.GitTagName,
		GitTagURL:    tagURL,
		AuthorUserID: authorUserID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := r.Validate(); err != nil {
		return Release{}, err
	}

	return r, nil
}

type UpdateReleaseFunc func(r Release) (Release, error)

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
	if r.GitTagName == "" {
		return errReleaseGitTagRequired
	}
	if r.GitTagURL == (url.URL{}) {
		return errReleaseGitTagURLRequired
	}

	return nil
}

type CreateReleaseAttachmentInput struct {
	Name     string
	FilePath string
}

func (i CreateReleaseAttachmentInput) Validate() error {
	if i.Name == "" {
		return errors.New("attachment name is required")
	}
	if i.FilePath == "" {
		return errors.New("attachment file path is required")
	}

	return nil
}

type ReleaseAttachment struct {
	ID        uuid.UUID
	Name      string
	FilePath  string
	URL       url.URL
	CreatedAt time.Time
}

func (r *ReleaseAttachment) Validate() error {
	if r.Name == "" {
		return errReleaseAttachmentNameRequired
	}
	if r.FilePath == "" {
		return errReleaseAttachmentFilePathRequired
	}
	if r.URL == (url.URL{}) {
		return errReleaseAttachmentURLRequired
	}

	return nil
}

func NewReleaseAttachment(input CreateReleaseAttachmentInput, fileURL url.URL) (ReleaseAttachment, error) {
	a := ReleaseAttachment{
		ID:        uuid.New(),
		Name:      input.Name,
		FilePath:  input.FilePath,
		URL:       fileURL,
		CreatedAt: time.Now(),
	}

	if err := a.Validate(); err != nil {
		return ReleaseAttachment{}, err
	}

	return a, nil
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
		n.GitTagName = &r.GitTagName
		n.GitTagURL = &r.GitTagURL
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

type GithubGeneratedReleaseNotesInput struct {
	GitTagName         *string
	PreviousGitTagName *string
}

func (i GithubGeneratedReleaseNotesInput) Validate() error {
	if i.GitTagName == nil {
		return errGithubGeneratedNotesGitTagRequired
	}
	if *i.GitTagName == "" {
		return errGithubGeneratedNotesGitTagRequired
	}

	return nil
}

func (i GithubGeneratedReleaseNotesInput) GetGitTagName() string {
	if i.GitTagName != nil {
		return *i.GitTagName
	}

	return ""
}

type GithubGeneratedReleaseNotes struct {
	Title string
	Notes string
}
