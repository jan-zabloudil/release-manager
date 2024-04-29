package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Resend struct {
	mock.Mock
}

func (r *Resend) SendProjectInvitationEmail(ctx context.Context, data model.ProjectInvitationInput) error {
	args := r.Called(ctx, data)
	return args.Error(0)
}
