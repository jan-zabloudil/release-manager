package util

import (
	"errors"
	"net/http"

	svcerrors "release-manager/service/errors"

	"github.com/google/go-github/v60/github"
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
