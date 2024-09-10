package storage

import (
	"fmt"
	"net/url"

	"github.com/nedpals/supabase-go"
)

const (
	// signedURLExpiration is the expiration time of the signed URL in seconds
	signedURLExpiration = 3600
)

type Client struct {
	client *supabase.Client
	bucket string
}

func NewClient(client *supabase.Client, bucket string) *Client {
	return &Client{
		client: client,
		bucket: bucket,
	}
}

func (c *Client) GenerateFileURL(filePath string) (url.URL, error) {
	resp, err := c.client.Storage.From(c.bucket).CreateSignedURLForDownload(filePath, signedURLExpiration)
	if err != nil {
		return url.URL{}, fmt.Errorf("creating signed URL for file %s: %w", filePath, err)
	}

	signedURL, err := url.Parse(resp.SignedURL)
	if err != nil {
		return url.URL{}, fmt.Errorf("parsing signed URL string: %w", err)
	}

	return *signedURL, nil
}
