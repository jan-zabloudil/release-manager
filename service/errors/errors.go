package errors

import "errors"

var (
	ErrInvalidUserRole     = errors.New("invalid user role")
	ErrInvalidTemplateType = errors.New("invalid template type")
)
