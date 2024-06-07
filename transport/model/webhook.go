package model

import (
	"errors"
	"strings"
)

const (
	repoSlugsCount = 2
)

var (
	errInvalidRepoSlugs = errors.New("invalid repository slugs")
)

type GithubRefWebhookInput struct {
	// Ref is the name of the reference (tag, branch, etc.)
	Ref string `json:"ref" required:"true"`
	// RefType is the type of the reference (e.g. "branch" or "tag")
	RefType string `json:"ref_type" required:"true"`
	Repo    struct {
		// Owner and repo slug of the GitHub repo separated by a slash
		// (e.g. "owner/repo")
		Slugs string `json:"full_name"`
	} `json:"repository"`
}

// SplitGithubRepoSlugs splits the owner and repo slugs from the full name of the GitHub repository.
// The full name should be in the format "owner/repo".
func SplitGithubRepoSlugs(slugs string) (ownerSlug, repoSlug string, err error) {
	s := strings.Split(slugs, "/")
	if len(s) != repoSlugsCount {
		return "", "", errInvalidRepoSlugs
	}

	ownerSlug = s[0]
	repoSlug = s[1]

	if ownerSlug == "" || repoSlug == "" {
		return "", "", errInvalidRepoSlugs
	}

	return ownerSlug, repoSlug, nil
}
