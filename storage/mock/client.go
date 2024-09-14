package mock

import (
	"net/url"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (m *Client) GenerateFileURL(filePath string) (url.URL, error) {
	args := m.Called(filePath)
	return args.Get(0).(url.URL), args.Error(1)
}

func (m *Client) FileExists(filePath string) (bool, error) {
	args := m.Called(filePath)
	return args.Bool(0), args.Error(1)
}
