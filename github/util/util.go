package util

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	svcerrors "release-manager/service/errors"

	"github.com/google/go-github/v60/github"
)

const (
	// GitHub API error codes
	errCodeAlreadyExists = "already_exists"

	// GitHub API error messages
	errMessageInvalidPreviousTag = "Invalid previous_tag parameter"

	// GitHub API error fields
	gitTagNameField = "tag_name"

	// expectedGithubRepositoryURLSlugCount is the expected number of slugs in a GitHub repository URL
	// Example URL: https://github.com/owner/repo -> owner and repo are the slugs
	expectedGithubRepositoryURLSlugCount = 2
)

var (
	errInvalidGithubRepoURLPath = errors.New("invalid GitHub repository URL path, not in the format /owner/repo")
)

func ParseGithubRepoURL(rawURL string) (ownerSlug, repoSlug string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	// GitHub repo URL format: https://github.com/owner/repo,
	// OwnerSlug: owner, RepoSlug: repo.
	path := strings.Trim(u.Path, "/")
	slugs := strings.Split(path, "/")

	if len(slugs) != expectedGithubRepositoryURLSlugCount {
		return "", "", errInvalidGithubRepoURLPath
	}

	if slugs[0] == "" || slugs[1] == "" {
		return "", "", errInvalidGithubRepoURLPath
	}

	return slugs[0], slugs[1], nil
}

// TranslateGithubAuthError translates GitHub auth errors to service errors
func TranslateGithubAuthError(err error) error {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) {
		switch githubErr.Response.StatusCode {
		case http.StatusUnauthorized:
			return svcerrors.NewGithubClientUnauthorizedError().Wrap(err)
		case http.StatusForbidden:
			return svcerrors.NewGithubClientForbiddenError().Wrap(err)
		}
	}

	return err
}

func IsNotFoundError(err error) bool {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) {
		return githubErr.Response.StatusCode == http.StatusNotFound
	}

	return false
}

func IsReleaseAlreadyExistsError(err error) bool {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) && githubErr.Errors != nil {
		// GitHub returns error response as an array of errors
		// Each error contains fields (code, resource, field)
		for _, e := range githubErr.Errors {
			if e.Code == errCodeAlreadyExists && e.Field == gitTagNameField {
				return true
			}
		}
	}

	return false
}

func IsInvalidPreviousTagError(err error) bool {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) {
		return githubErr.Message == errMessageInvalidPreviousTag
	}

	return false
}

func GenerateRepoURL(ownerSlug, repoSlug string) (url.URL, error) {
	if ownerSlug == "" || repoSlug == "" {
		return url.URL{}, errors.New("empty owner or repo slug")
	}

	rawURL := fmt.Sprintf("https://github.com/%s/%s", ownerSlug, repoSlug)
	u, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, err
	}

	return *u, nil
}

func GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error) {
	if tagName == "" || ownerSlug == "" || repoSlug == "" {
		return url.URL{}, errors.New("empty tag name, owner or repo slug")
	}

	// For ReleaseManager's users it is the most beneficial to see GitHub tag page (that is also a release page)
	// This page is available even if GitHub release is not created yet
	rawURL := fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", ownerSlug, repoSlug, tagName)
	u, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, err
	}

	return *u, nil
}
