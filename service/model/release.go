package model

import (
	"errors"
	"net/url"
)

var (
	errTagNameRequired               = errors.New("tag name is required")
	errReleaseNameRequired           = errors.New("release name is required")
	errReleaseSourceCodeLinkRequired = errors.New("release must be linked to a source code")
	errReleaseSourceCodeLinkConflict = errors.New("release can be linked to either a new git tag or an existing git tag, not both")
	errGithubReleaseNameRequired     = errors.New("github release name is required")
	errGithubReleasePageURLRequired  = errors.New("github release page URL is required")
)

type CreateReleaseDraftInput struct {
	TagName string
	// The commitish value, which can be a branch name or a commit SHA, specifies the base from which a Git tag is created (if it does not exist yet)
	// If not provided, the default branch of the repository is used
	TargetCommitish *string
	PreviousTagName *string
}

func (c CreateReleaseDraftInput) Validate() error {
	if c.TagName == "" {
		return errTagNameRequired
	}

	return nil
}

type GitTag struct {
	Name string
}

type GitTagInput struct {
	TagName         string
	TargetCommitish string
}

type GithubRelease struct {
	Name           string
	ReleaseNotes   string
	ReleasePageURL url.URL // Link to the release page on GitHub
}

func (g GithubRelease) Validate() error {
	if g.Name == "" {
		return errGithubReleaseNameRequired
	}

	if g.ReleasePageURL == (url.URL{}) {
		return errGithubReleasePageURLRequired
	}

	return nil
}

type ReleaseDraft struct {
	Name           string
	ReleaseNotes   string
	SourceCodeLink struct { // links the release with the source code, either by a new git tag or an existing git tag
		NewGitTag      *GitTagInput
		ExistingGitTag *GitTag
	}
	GithubRelease *GithubRelease
}

func (r *ReleaseDraft) Validate() error {
	if r.Name == "" {
		return errReleaseNameRequired
	}

	if r.SourceCodeLink.NewGitTag == nil && r.SourceCodeLink.ExistingGitTag == nil {
		return errReleaseSourceCodeLinkRequired
	}

	if r.SourceCodeLink.NewGitTag != nil && r.SourceCodeLink.ExistingGitTag != nil {
		return errReleaseSourceCodeLinkConflict
	}

	if r.GithubRelease != nil {
		if err := r.GithubRelease.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func NewDraftRelease(name, notes string) ReleaseDraft {
	return ReleaseDraft{
		Name:         name,
		ReleaseNotes: notes,
	}
}

func (r *ReleaseDraft) LinkSourceCodeByNewTag(g GitTagInput) {
	r.SourceCodeLink.NewGitTag = &g
}

func (r *ReleaseDraft) LinkSourceCodeByExistingTag(t GitTag) {
	r.SourceCodeLink.ExistingGitTag = &t
}

func (r *ReleaseDraft) AddGithubRelease(gr GithubRelease) {
	r.GithubRelease = &gr
}
