# Plan: Implement `delete_time_entry` MCP Tool

## Overview
Add the `delete_time_entry` tool to delete a time entry by ID.

## API Details
- **Endpoint:** `DELETE /workspaces/{workspaceId}/time-entries/{id}`
- **Response:** 204 No Content on success

## Files to Modify

### 1. `clockify/client.go`
Add DELETE method and DeleteTimeEntry:
```go
func (c *Client) Delete(ctx context.Context, path string) error

func (c *Client) DeleteTimeEntry(ctx context.Context, workspaceID, entryID string) error
```
- DELETE method doesn't need to decode response body
- Returns nil on success (204), error otherwise

### 2. `clockify/tools.go`
- `DeleteTimeEntryTool()` - tool definition with parameters:
  - `entry_id` (required) - ID of the time entry to delete
- `HandleDeleteTimeEntry()` - handler:
  1. Fetch current user for workspace ID
  2. Call `DeleteTimeEntry()`
  3. Return confirmation message

### 3. `main.go`
Register: `s.AddTool(clockify.DeleteTimeEntryTool(), tools.HandleDeleteTimeEntry)`

## Verification
1. Create a test time entry
2. Use `list_time_entries` to get its ID
3. Delete it with `delete_time_entry`
4. Verify it's gone from `list_time_entries`
