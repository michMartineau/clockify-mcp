package clockify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// Post performs a POST request to the given path with the body and decodes the response into v.
func (c *Client) Post(ctx context.Context, path string, body any, v any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, bytes.NewReader(jsonBody))
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
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// Patch performs a PATCH request to the given path with the body and decodes the response into v.
func (c *Client) Patch(ctx context.Context, path string, body any, v any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, baseURL+path, bytes.NewReader(jsonBody))
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
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
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

// GetProjects fetches all projects in a workspace.
func (c *Client) GetProjects(ctx context.Context, workspaceID string) ([]Project, error) {
	var projects []Project
	path := fmt.Sprintf("/workspaces/%s/projects", workspaceID)
	if err := c.Get(ctx, path, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// GetTags fetches all tags in a workspace.
func (c *Client) GetTags(ctx context.Context, workspaceID string) ([]Tag, error) {
	var tags []Tag
	path := fmt.Sprintf("/workspaces/%s/tags", workspaceID)
	if err := c.Get(ctx, path, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// TimeEntriesParams holds the query parameters for fetching time entries.
type TimeEntriesParams struct {
	Start     string // ISO 8601 datetime
	End       string // ISO 8601 datetime
	ProjectID string // Optional project filter
	Page      int
	PageSize  int
}

// GetTimeEntries fetches time entries for a user in a workspace with optional filters.
func (c *Client) GetTimeEntries(ctx context.Context, workspaceID, userID string, params TimeEntriesParams) ([]TimeEntry, error) {
	basePath := fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID)

	// Build query parameters
	query := url.Values{}
	if params.Start != "" {
		query.Set("start", params.Start)
	}
	if params.End != "" {
		query.Set("end", params.End)
	}
	if params.ProjectID != "" {
		query.Set("project", params.ProjectID)
	}
	if params.Page > 0 {
		query.Set("page", fmt.Sprintf("%d", params.Page))
	}
	if params.PageSize > 0 {
		query.Set("page-size", fmt.Sprintf("%d", params.PageSize))
	}

	path := basePath
	if len(query) > 0 {
		path = basePath + "?" + query.Encode()
	}

	var entries []TimeEntry
	if err := c.Get(ctx, path, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// GetCurrentTimer fetches the currently running time entry, if any.
// Returns nil if no timer is running.
func (c *Client) GetCurrentTimer(ctx context.Context, workspaceID, userID string) (*TimeEntry, error) {
	var entries []TimeEntry
	path := fmt.Sprintf("/workspaces/%s/user/%s/time-entries?in-progress=true", workspaceID, userID)
	if err := c.Get(ctx, path, &entries); err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	return &entries[0], nil
}

// StartTimer starts a new time entry (running timer).
func (c *Client) StartTimer(ctx context.Context, workspaceID string, req StartTimerRequest) (*TimeEntry, error) {
	var entry TimeEntry
	path := fmt.Sprintf("/workspaces/%s/time-entries", workspaceID)
	if err := c.Post(ctx, path, req, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// StopTimer stops the currently running time entry.
func (c *Client) StopTimer(ctx context.Context, workspaceID, userID string, endTime string) (*TimeEntry, error) {
	var entry TimeEntry
	path := fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID)
	body := map[string]string{"end": endTime}
	if err := c.Patch(ctx, path, body, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Delete performs a DELETE request to the given path.
// Returns nil on success (expects 204 No Content).
func (c *Client) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL+path, nil)
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

	return nil
}

// CreateTimeEntry creates a completed time entry with both start and end times.
func (c *Client) CreateTimeEntry(ctx context.Context, workspaceID string, req CreateTimeEntryRequest) (*TimeEntry, error) {
	var entry TimeEntry
	path := fmt.Sprintf("/workspaces/%s/time-entries", workspaceID)
	if err := c.Post(ctx, path, req, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// DeleteTimeEntry deletes a time entry by ID.
func (c *Client) DeleteTimeEntry(ctx context.Context, workspaceID, entryID string) error {
	path := fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, entryID)
	return c.Delete(ctx, path)
}