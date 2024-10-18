package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"release-manager/pkg/id"

	"github.com/nedpals/supabase-go"
)

var (
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
)

type Client struct {
	client *supabase.Client
}

func NewClient(client *supabase.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Authenticate(ctx context.Context, token string) (id.AuthUser, error) {
	user, err := c.client.Auth.User(ctx, token)
	if err != nil {
		var errResponse *supabase.ErrorResponse
		if errors.As(err, &errResponse) && errResponse.Code == http.StatusForbidden {
			return id.AuthUser{}, fmt.Errorf("%w: %s", ErrInvalidOrExpiredToken, errResponse.Message)
		}

		return id.AuthUser{}, fmt.Errorf("authenticating user: %w", err)
	}

	var userID id.AuthUser
	if err := userID.FromString(user.ID); err != nil {
		return id.AuthUser{}, fmt.Errorf("parsing user ID: %w", err)
	}

	return userID, nil
}
