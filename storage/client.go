package storage

import (
	"fmt"
	"net/url"

	"release-manager/storage/util"

	"github.com/nedpals/supabase-go"
)

const (
	// signedURLExpiration is the expiration time of the signed URL in seconds
	signedURLExpiration = 3600
)

var (
	errInvalidFileKey = fmt.Errorf("file key must be in the format bucketID/filePath")
)

type Client struct {
	client *supabase.Client
}

func NewClient(client *supabase.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) GenerateFileURL(fileKey string) (url.URL, error) {
	bucketID, filePath := util.ExplodeFileKey(fileKey)
	if bucketID == "" || filePath == "" {
		return url.URL{}, fmt.Errorf("parsing file key %s: %w", fileKey, errInvalidFileKey)
	}

	resp, err := c.client.Storage.From(bucketID).CreateSignedURLForDownload(filePath, signedURLExpiration)
	if err != nil {
		return url.URL{}, fmt.Errorf("creating signed URL for file %s: %w", fileKey, err)
	}

	signedURL, err := url.Parse(resp.SignedURL)
	if err != nil {
		return url.URL{}, fmt.Errorf("parsing signed URL string: %w", err)
	}

	return *signedURL, nil
}
