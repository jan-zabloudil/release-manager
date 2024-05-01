package util

import (
	"errors"
	"net/http"

	"release-manager/pkg/githuberrors"

	"github.com/google/go-github/v60/github"
)

func ToGithubError(err error) *githuberrors.GithubError {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) {
		switch githubErr.Response.StatusCode {
		case http.StatusUnauthorized:
			return githuberrors.NewUnauthorizedError().Wrap(err)
		case http.StatusForbidden:
			return githuberrors.NewForbiddenError().Wrap(err)
		case http.StatusNotFound:
			return githuberrors.NewNotFoundError().Wrap(err)
		}
	}

	return githuberrors.NewUnknownError().Wrap(err)
}
