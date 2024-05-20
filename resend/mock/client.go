package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (c *Client) SendProjectInvitationEmailAsync(ctx context.Context, data model.ProjectInvitationEmailData, recipient string) {
	c.Called(ctx, data, recipient)
}
