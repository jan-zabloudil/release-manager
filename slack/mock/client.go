package slack

import (
	"context"

	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (m *Client) SendReleaseNotification(ctx context.Context, tkn model.SlackToken, channelID string, n model.ReleaseNotification) error {
	args := m.Called(ctx, tkn, channelID, n)
	return args.Error(0)
}
