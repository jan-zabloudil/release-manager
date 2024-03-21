package errors

import "errors"

var (
	UnknownErr                                        = errors.New("supabase: unknown error")
	ErrEntityNotExists                                = errors.New("supabase: entity does not exist")
	ErrUserAuthenticationFailed                       = errors.New("supabase: user authentication by bearer token unsuccessful")
	ErrAccessToResourceDenied                         = errors.New("supabase: access to resource denied")
	ErrMultipleOrNoRecordsReturnedAfterWriteOperation = errors.New("supabase: one record expected after insert/update operation, multiple (or no) rows returned")
	ErrResourceNotFound                               = errors.New("supabase: resource not found ")
)
