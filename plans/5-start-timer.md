# Plan: Implement `start_timer` MCP Tool

## Overview
Add the `start_timer` tool to start a new time entry (running timer).

## API Details
- **Endpoint:** `POST /workspaces/{workspaceId}/time-entries`
- **Request Body:**
```json
{
  "start": "2024-01-15T09:00:00Z",
  "description": "Working on feature",
  "projectId": "optional-project-id",
  "tagIds": ["optional-tag-id"],
  "billable": true
}
```
- **Note:** Omit `end` field to create a running timer

## Files to Modify

### 1. `clockify/client.go`
Add POST method and StartTimer:
```go
func (c *Client) Post(ctx context.Context, path string, body any, v any) error

func (c *Client) StartTimer(ctx context.Context, workspaceID string, req StartTimerRequest) (*TimeEntry, error)
```

Add request struct:
```go
type StartTimerRequest struct {
    Start       string   `json:"start"`
    Description string   `json:"description,omitempty"`
    ProjectID   string   `json:"projectId,omitempty"`
    TagIDs      []string `json:"tagIds,omitempty"`
    Billable    bool     `json:"billable,omitempty"`
}
```

### 2. `clockify/tools.go`
- `StartTimerTool()` - tool definition with parameters:
  - `description` (optional) - what you're working on
  - `project_id` (optional) - project to associate
  - `tag_ids` (optional) - array of tag IDs
  - `billable` (optional) - boolean
- `HandleStartTimer()` - handler:
  1. Fetch current user for default workspace
  2. Build request with current time as start
  3. Call `StartTimer()`
  4. Return confirmation with entry details

### 3. `main.go`
Register: `s.AddTool(clockify.StartTimerTool(), tools.HandleStartTimer)`

## Verification
1. Build and test starting timer with just description
2. Verify timer appears in Clockify UI
3. Test with project_id and tags