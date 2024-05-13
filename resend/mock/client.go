package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (c *Client) SendProjectInvitationEmail(ctx context.Context, data model.ProjectInvitationInput) error {
	args := c.Called(ctx, data)
	return args.Error(0)
}
