# Plan: Implement `list_time_entries` MCP Tool

## Overview
Add the `list_time_entries` tool to list time entries for a date range.

## API Details
- **Endpoint:** `GET /workspaces/{workspaceId}/user/{userId}/time-entries`
- **Query Parameters:**
  - `start` - ISO 8601 datetime (required by tool, defaults to start of today)
  - `end` - ISO 8601 datetime (optional)
  - `project` - filter by project ID (optional)
  - `page`, `page-size` - pagination

## Files to Modify

### 1. `clockify/client.go`
Add method with query params support:
```go
func (c *Client) GetTimeEntries(ctx context.Context, workspaceID, userID string, params TimeEntriesParams) ([]TimeEntry, error)

type TimeEntriesParams struct {
    Start     string // ISO 8601
    End       string // ISO 8601
    ProjectID string
    Page      int
    PageSize  int
}
```

May need to extend `Get` method or add `GetWithParams` to handle query strings.

### 2. `clockify/tools.go`
- `ListTimeEntriesTool()` - tool definition with parameters:
  - `start_date` (optional) - defaults to start of today
  - `end_date` (optional) - defaults to now
  - `project_id` (optional) - filter by project
- `HandleListTimeEntries()` - handler:
  1. Fetch current user for workspace/user IDs
  2. Parse/default date parameters
  3. Call `GetTimeEntries()`
  4. Format results showing description, project, duration, times

### 3. `main.go`
Register: `s.AddTool(clockify.ListTimeEntriesTool(), tools.HandleListTimeEntries)`

## Verification
1. Build and test with no params (should return today's entries)
2. Test with date range
3. Test with project filter
4. Verify pagination works for users with many entries