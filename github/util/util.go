package util

import (
	"errors"
	"net/http"

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
)

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

func IsGithubNotFoundError(err error) bool {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) {
		return githubErr.Response.StatusCode == http.StatusNotFound
	}

	return false
}

func IsGithubReleaseAlreadyExistsError(err error) bool {
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

func IsGithubInvalidPreviousTagError(err error) bool {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) {
		return githubErr.Message == errMessageInvalidPreviousTag
	}

	return false
}
