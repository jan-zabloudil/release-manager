package model

import (
	"fmt"

	cryptox "release-manager/pkg/crypto"
)

type Email struct {
	Subject    string
	Text       string
	HTML       string
	Recipients []string
}

func NewProjectInvitationEmail(p Project, tkn cryptox.Token, recipient string) Email {
	// TODO implement proper email template (both html and text) with proper message and magic link, this is just a starting point!
	return Email{
		Subject:    "Project invitation",
		Text:       fmt.Sprintf("You have been invited to join the project %s. (Token: %s)", p.Name, tkn),
		HTML:       fmt.Sprintf("<p>You have been invited to join the project %s.</p><p>(Token: %s)</p>", p.Name, tkn),
		Recipients: []string{recipient},
	}
}
