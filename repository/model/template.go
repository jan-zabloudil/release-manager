package model

import (
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Template struct {
	ID         uuid.UUID      `json:"id"`
	Type       string         `json:"type"`
	ReleaseMsg ReleaseMessage `json:"template_data"`
}

func ToDBTemplate(id uuid.UUID, tmplType string, msg ReleaseMessage) Template {
	return Template{
		ID:         id,
		Type:       tmplType,
		ReleaseMsg: msg,
	}
}

func ToSvcTemplate(id uuid.UUID, typeStr string, msg svcmodel.ReleaseMessage) (svcmodel.Template, error) {
	tmplType, err := svcmodel.NewTemplateType(typeStr)
	if err != nil {
		return svcmodel.Template{}, err
	}

	return svcmodel.Template{
		ID:         id,
		Type:       tmplType,
		ReleaseMsg: msg,
	}, nil
}

func ToSvcTemplates(templates []Template) ([]svcmodel.Template, error) {
	t := make([]svcmodel.Template, 0, len(templates))
	for _, template := range templates {
		svcTmpl, err := ToSvcTemplate(
			template.ID,
			template.Type,
			ToSvcReleaseMsg(template.ReleaseMsg.Title, template.ReleaseMsg.Text, template.ReleaseMsg.Includes),
		)
		if err != nil {
			return nil, err
		}

		t = append(t, svcTmpl)
	}

	return t, nil
}
