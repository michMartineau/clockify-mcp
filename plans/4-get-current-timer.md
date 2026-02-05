# Plan: Implement `get_current_timer` MCP Tool

## Overview
Add the `get_current_timer` tool to retrieve the currently running time entry.

## API Details
- **Endpoint:** `GET /workspaces/{workspaceId}/user/{userId}/time-entries?in-progress=true`
- **Response:** Array of time entries (0 or 1 element when in-progress=true)

## Files to Modify

### 1. `clockify/models.go`
Add `TimeEntry` struct:
```go
type TimeEntry struct {
    ID           string     `json:"id"`
    Description  string     `json:"description"`
    ProjectID    string     `json:"projectId"`
    TagIDs       []string   `json:"tagIds"`
    Billable     bool       `json:"billable"`
    TimeInterval TimeInterval `json:"timeInterval"`
}

type TimeInterval struct {
    Start    string `json:"start"`    // ISO 8601
    End      string `json:"end"`      // ISO 8601, empty if running
    Duration string `json:"duration"` // ISO 8601 duration
}
```

### 2. `clockify/client.go`
Add method:
```go
func (c *Client) GetCurrentTimer(ctx context.Context, workspaceID, userID string) (*TimeEntry, error)
```
- Returns `nil` if no timer running (empty array from API)

### 3. `clockify/tools.go`
- `GetCurrentTimerTool()` - tool definition (no required params)
- `HandleGetCurrentTimer()` - handler:
  1. Fetch current user for `workspaceId` and `userId`
  2. Call `GetCurrentTimer()`
  3. Return timer info or "No timer running"

### 4. `main.go`
Register: `s.AddTool(clockify.GetCurrentTimerTool(), tools.HandleGetCurrentTimer)`

## Verification
1. Build and test with no timer running (should say "No timer running")
2. Start a timer in Clockify UI, verify tool returns it