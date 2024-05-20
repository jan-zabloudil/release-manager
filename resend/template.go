package resend

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	"release-manager/config"
	svcmodel "release-manager/service/model"
)

//go:embed templates/project_invitation.tmpl
var projectInvitationTmpl string

type ParsedTemplate struct {
	Subject string
	Text    string
	HTML    string
}

func ParseProjectInvitationTemplate(data svcmodel.ProjectInvitationEmailData, clientSvcCfg config.ClientServiceConfig) (ParsedTemplate, error) {
	tmpl, err := template.New("project_invitation").Parse(projectInvitationTmpl)
	if err != nil {
		return ParsedTemplate{}, fmt.Errorf("failed to parse templates templates: %w", err)
	}

	templateData := map[string]string{
		"projectName": data.ProjectName,
		"siteLink":    clientSvcCfg.URL,
		"acceptLink":  fmt.Sprintf("%s/%s?token=%s", clientSvcCfg.URL, clientSvcCfg.AcceptInvitationRoute, data.Token),
		"rejectLink":  fmt.Sprintf("%s/%s?token=%s", clientSvcCfg.URL, clientSvcCfg.RejectInvitationRoute, data.Token),
		"signUpLink":  fmt.Sprintf("%s/%s", clientSvcCfg.URL, clientSvcCfg.SignUpRoute),
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", templateData)
	if err != nil {
		return ParsedTemplate{}, fmt.Errorf("failed to execute templates templates: %w", err)
	}

	textBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(textBody, "textBody", templateData)
	if err != nil {
		return ParsedTemplate{}, fmt.Errorf("failed to execute templates templates: %w", err)
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", templateData)
	if err != nil {
		return ParsedTemplate{}, fmt.Errorf("failed to execute templates templates: %w", err)
	}

	return ParsedTemplate{
		Subject: subject.String(),
		Text:    textBody.String(),
		HTML:    htmlBody.String(),
	}, nil
}
