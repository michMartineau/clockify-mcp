# Plan: Implement `stop_timer` MCP Tool

## Overview
Add the `stop_timer` tool to stop the currently running time entry.

## API Details
- **Endpoint:** `PATCH /workspaces/{workspaceId}/user/{userId}/time-entries`
- **Request Body:**
```json
{
  "end": "2024-01-15T17:00:00Z"
}
```

## Files to Modify

### 1. `clockify/client.go`
Add PATCH method and StopTimer:
```go
func (c *Client) Patch(ctx context.Context, path string, body any, v any) error

func (c *Client) StopTimer(ctx context.Context, workspaceID, userID string, endTime string) (*TimeEntry, error)
```

### 2. `clockify/tools.go`
- `StopTimerTool()` - tool definition (no required params)
- `HandleStopTimer()` - handler:
  1. Fetch current user for workspace/user IDs
  2. Call `StopTimer()` with current time as end
  3. Return stopped entry details or error if no timer running

### 3. `main.go`
Register: `s.AddTool(clockify.StopTimerTool(), tools.HandleStopTimer)`

## Verification
1. Start a timer, then test stopping it
2. Verify stopped entry in Clockify UI has correct end time
3. Test stopping when no timer running (should return clear error)