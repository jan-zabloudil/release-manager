package errors

import "errors"

var (
	ErrUserIsAlreadyMember              = errors.New("user with provided email is already a project member")
	ErrInvitationAlreadyExists          = errors.New("project invitation for this email already exists")
	ErrInvalidProjectRole               = errors.New("invalid project role")
	ErrProjectMemberRoleCannotBeGranted = errors.New("project member does not have permission to assign this role to others")
	ErrProjectMemberUpdateNotAllowed    = errors.New("project member cannot update another member unless the other member has an inferior role")
	ErrAppEnvURLInvalid                 = errors.New("app environment URL is not valid absolute url")
)
