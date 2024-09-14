package storage

import (
	"errors"
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

func (c *Client) FileExists(filePath string) (bool, error) {
	// Supabase API does not have a direct method to check if a file exists
	// The lightest way to check if a file exists is to get the file metadata
	// If file does not exist, it will return ErrNotFound error
	_, err := c.client.Storage.From(c.bucket).GetFileMetadata(filePath)
	if err != nil {
		if errors.Is(err, supabase.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("getting metadata for file %s: %w", filePath, err)
	}

	return true, nil
}
