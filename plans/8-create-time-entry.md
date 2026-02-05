# Plan: Implement `create_time_entry` MCP Tool

## Overview
Add the `create_time_entry` tool to log a completed time entry with both start and end times.

## API Details
- **Endpoint:** `POST /workspaces/{workspaceId}/time-entries`
- **Request Body:**
```json
{
  "start": "2024-01-15T09:00:00Z",
  "end": "2024-01-15T17:00:00Z",
  "description": "Worked on feature",
  "projectId": "optional-project-id",
  "tagIds": ["optional-tag-id"],
  "billable": true
}
```
- **Note:** Including `end` field creates a completed entry (vs running timer)

## Files to Modify

### 1. `clockify/models.go`
Add request struct (or reuse/extend `StartTimerRequest`):
```go
type CreateTimeEntryRequest struct {
    Start       string   `json:"start"`
    End         string   `json:"end"`
    Description string   `json:"description,omitempty"`
    ProjectID   string   `json:"projectId,omitempty"`
    TagIDs      []string `json:"tagIds,omitempty"`
    Billable    bool     `json:"billable,omitempty"`
}
```

### 2. `clockify/client.go`
Add method:
```go
func (c *Client) CreateTimeEntry(ctx context.Context, workspaceID string, req CreateTimeEntryRequest) (*TimeEntry, error)
```
- Uses existing `Post()` method

### 3. `clockify/tools.go`
- `CreateTimeEntryTool()` - tool definition with parameters:
  - `start` (required) - ISO 8601 start time
  - `end` (required) - ISO 8601 end time
  - `description` (optional)
  - `project_id` (optional)
  - `billable` (optional)
- `HandleCreateTimeEntry()` - validates start < end, calls client, returns entry details

### 4. `main.go`
Register: `s.AddTool(clockify.CreateTimeEntryTool(), tools.HandleCreateTimeEntry)`

## Verification
1. Build and test creating an entry for yesterday
2. Verify entry appears in Clockify UI with correct times
3. Test validation rejects end before start
