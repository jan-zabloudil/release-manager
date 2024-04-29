package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type EmailService struct {
	mock.Mock
}

func (m *EmailService) SendProjectInvitation(ctx context.Context, input model.ProjectInvitationInput) {
	m.Called(ctx, input)
}
