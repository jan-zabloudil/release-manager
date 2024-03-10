package utils

import (
	"errors"
	"fmt"
	"net/http"

	githuberr "release-manager/github/errors"

	"github.com/google/go-github/v60/github"
)

func WrapGithubErr(err error) error {
	var errResp *github.ErrorResponse
	if errors.As(err, &errResp) {
		for _, err := range errResp.Errors {
			switch err.Code {
			case "already_exists":
				return fmt.Errorf("%w: %s", githuberr.ErrTagAlreadyExists, err.Error())
			case "invalid":
				return fmt.Errorf("%w: %s", githuberr.ErrInvalidTargetCommitish, err.Error())
			}
		}

		switch errResp.Response.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("%w: %s", githuberr.ErrResourceNotFound, err.Error())
		case http.StatusUnauthorized:
			return fmt.Errorf("%w: %s", githuberr.ErrUnauthenticated, err.Error())
		case http.StatusForbidden:
			return fmt.Errorf("%w: %s", githuberr.ErrForbidden, err.Error())
		}
	}

	return fmt.Errorf("%w: %s", githuberr.UnknownErr, err.Error())
}
