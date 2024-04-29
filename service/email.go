package service

import (
	"context"
	"fmt"

	"release-manager/service/model"
)

type EmailService struct {
	mailer mailer
}

func NewEmailService(m mailer) *EmailService {
	return &EmailService{mailer: m}
}

func (s *EmailService) SendProjectInvitation(ctx context.Context, input model.ProjectInvitationInput) {
	// TODO implement proper email template (both html and text) with proper message and magic link, this is just a starting point!
	subject := "Project invitation"
	text := fmt.Sprintf("You have been invited to join the project %s. (Token: %s)", input.ProjectName, input.Token)
	html := fmt.Sprintf("<p>You have been invited to join the project %s.</p><p>Token: %s)</p>", input.ProjectName, input.Token)

	s.mailer.SendEmailAsync(ctx, subject, text, html, input.RecipientEmail)
}
