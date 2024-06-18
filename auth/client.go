package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func (c *Client) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	user, err := c.client.Auth.User(ctx, token)
	if err != nil {
		var errResponse *supabase.ErrorResponse
		if errors.As(err, &errResponse) && errResponse.Code == http.StatusForbidden {
			return uuid.UUID{}, fmt.Errorf("%w: %s", ErrInvalidOrExpiredToken, errResponse.Message)
		}

		return uuid.UUID{}, fmt.Errorf("authenticating user: %w", err)
	}

	id, err := uuid.Parse(user.ID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parsing user ID: %w", err)
	}

	return id, nil
}
