package slack

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (m *Client) SendReleaseNotification(ctx context.Context, token, channelID string, n model.ReleaseNotification) {
	m.Called(ctx, token, channelID, n)
}
