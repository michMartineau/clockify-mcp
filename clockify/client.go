package clockify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "https://api.clockify.me/api/v1"

// Client is an HTTP client for the Clockify API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Clockify client.
// It reads the API key from the CLOCKIFY_API_KEY environment variable.
func NewClient() (*Client, error) {
	apiKey := os.Getenv("CLOCKIFY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CLOCKIFY_API_KEY environment variable is not set")
	}
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

// Get performs a GET request to the given path and decodes the response into v.
func (c *Client) Get(ctx context.Context, path string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// GetCurrentUser fetches the current user's information.
func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var user User
	if err := c.Get(ctx, "/user", &user); err != nil {
		return nil, err
	}
	return &user, nil
}