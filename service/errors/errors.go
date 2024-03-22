package errors

import "errors"

var (
	ErrUserIsAlreadyMember              = errors.New("user with provided email is already a project member")
	ErrInvitationAlreadyExists          = errors.New("project invitation for this email already exists")
	ErrInvalidProjectRole               = errors.New("invalid project role")
	ErrProjectMemberRoleCannotBeGranted = errors.New("project member does not have permission to assign this role to others")
	ErrProjectMemberUpdateNotAllowed    = errors.New("project member cannot update another member unless the other member has an inferior role")
	ErrAppEnvURLInvalid                 = errors.New("app environment URL is not valid absolute url")
	ErrUnknownSCMRepoPlatform           = errors.New("unknown SCM repo platform")
	ErrInvalidSCMRepoURL                = errors.New("invalid SCM repo URL")
	ErrInvalidGithubHostUrl             = errors.New("repo url must have github.com as host")
	ErrInvalidGithubRepoUrl             = errors.New("invalid github repository url")
	ErrSCMRepoNotSet                    = errors.New("SCM repo is not set")
	ErrInvalidTag                       = errors.New("invalid tag name")
)
