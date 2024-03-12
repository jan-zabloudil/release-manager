package errors

import "errors"

var (
	ErrUserIsAlreadyMember     = errors.New("user with provided email is already a project member")
	ErrInvitationAlreadyExists = errors.New("project invitation for this email already exists")
	ErrInvalidProjectRole      = errors.New("invalid project role")
)
