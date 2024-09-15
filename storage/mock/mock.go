package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (m *Client) DeleteFileAsync(ctx context.Context, filePath string) {
	m.Called(ctx, filePath)
}
