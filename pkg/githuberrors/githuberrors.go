package githuberrors

import (
	"errors"
	"fmt"
)

const (
	errCodeUnauthorized        = "ERR_GITHUB_UNAUTHORIZED"
	errCodeForbidden           = "ERR_GITHUB_FORBIDDEN"
	errCodeNotFound            = "ERR_GITHUB_NOT_FOUND"
	errCodeCannotMapToSvcModel = "ERR_GITHUB_CANNOT_MAP_TO_SVC_MODEL"
	errCodeUnknown             = "ERR_GITHUB_UNKNOWN"
)

type GithubError struct {
	Code string
	Err  error
}

func (e *GithubError) Error() string {
	return fmt.Sprintf("Code: %s, error: %s", e.Code, e.Err)
}

func (e *GithubError) Wrap(err error) *GithubError {
	return &GithubError{
		Code: e.Code,
		Err:  err,
	}
}

func NewUnauthorizedError() *GithubError {
	return &GithubError{
		Code: errCodeUnauthorized,
	}
}

func NewForbiddenError() *GithubError {
	return &GithubError{
		Code: errCodeForbidden,
	}
}

func NewNotFoundError() *GithubError {
	return &GithubError{
		Code: errCodeNotFound,
	}
}

func NewUnknownError() *GithubError {
	return &GithubError{
		Code: errCodeUnknown,
	}
}

func NewToSvcModelError() *GithubError {
	return &GithubError{
		Code: errCodeCannotMapToSvcModel,
	}
}

func IsErrorWithCode(err error, code string) bool {
	var githubErr *GithubError
	if errors.As(err, &githubErr) {
		return githubErr.Code == code
	}

	return false
}

func IsUnauthorizedError(err error) bool {
	return IsErrorWithCode(err, errCodeUnauthorized)
}

func IsForbiddenError(err error) bool {
	return IsErrorWithCode(err, errCodeForbidden)
}

func IsNotFoundError(err error) bool {
	return IsErrorWithCode(err, errCodeNotFound)
}
