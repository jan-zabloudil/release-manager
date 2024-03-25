package model

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type TemplateService interface {
	Create(ctx context.Context, t svcmodel.Template) (svcmodel.Template, error)
	Get(ctx context.Context, id uuid.UUID) (svcmodel.Template, error)
	GetAll(ctx context.Context) ([]svcmodel.Template, error)
	Update(ctx context.Context, t svcmodel.Template) (svcmodel.Template, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Template struct {
	ID         uuid.UUID      `json:"id"`
	Type       string         `json:"type" validate:"required"`
	ReleaseMsg ReleaseMessage `json:"template_data"`
}

func NewSvcTemplate(typeStr string, msg svcmodel.ReleaseMessage) (svcmodel.Template, error) {
	tmplType, err := svcmodel.NewTemplateType(typeStr)
	if err != nil {
		return svcmodel.Template{}, err
	}

	t := svcmodel.Template{
		ID:   uuid.New(),
		Type: tmplType,
	}

	return ToSvcTemplate(t, msg), nil
}

func ToSvcTemplate(t svcmodel.Template, msg svcmodel.ReleaseMessage) svcmodel.Template {
	t.ReleaseMsg = msg

	return t
}

func ToNetTemplate(id uuid.UUID, tmplType string, msg ReleaseMessage) Template {
	return Template{
		ID:         id,
		Type:       tmplType,
		ReleaseMsg: msg,
	}
}

func ToNetTemplates(templates []svcmodel.Template) []Template {
	t := make([]Template, 0, len(templates))
	for _, template := range templates {
		t = append(t, ToNetTemplate(
			template.ID,
			template.Type.TemplateType(),
			ToNetReleaseMsg(template.ReleaseMsg.Title, template.ReleaseMsg.Text, template.ReleaseMsg.Includes)),
		)
	}

	return t
}
