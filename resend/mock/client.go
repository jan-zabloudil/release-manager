package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (c *Client) SendEmailAsync(ctx context.Context, email model.Email) {
	c.Called(ctx, email)
}
