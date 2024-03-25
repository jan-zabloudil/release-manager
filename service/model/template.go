package model

import (
	"context"

	svcerr "release-manager/service/errors"

	"github.com/google/uuid"
)

const (
	releaseMsgTmplType = "release_message"
)

type TemplateRepository interface {
	Insert(ctx context.Context, t Template) (Template, error)
	ReadAll(ctx context.Context) ([]Template, error)
	Read(ctx context.Context, id uuid.UUID) (Template, error)
	Update(ctx context.Context, t Template) (Template, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Template struct {
	ID         uuid.UUID
	Type       TemplateType
	ReleaseMsg ReleaseMessage
}

type TemplateType interface {
	TemplateType() string
}

type templateType struct {
	templateType string
}

func (t templateType) TemplateType() string {
	return t.templateType
}

func NewTemplateType(key string) (TemplateType, error) {
	switch {
	case key == releaseMsgTmplType:
		return &templateType{templateType: releaseMsgTmplType}, nil
	default:
		return nil, svcerr.ErrInvalidTemplateType
	}
}
