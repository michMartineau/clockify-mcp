# Plan: Implement `get_weekly_summary` MCP Tool

## Overview
Add the `get_weekly_summary` tool to aggregate time by project for the current week.

## API Details
- **Uses:** `GET /workspaces/{workspaceId}/user/{userId}/time-entries?start=...&end=...`
- **Logic:** Fetch all entries for the week (Monday-Sunday), aggregate by project and day

## Files to Modify

### 1. `clockify/tools.go`
- `GetWeeklySummaryTool()` - tool definition with parameters:
  - `week_offset` (optional) - 0 for current week, -1 for last week, etc. Default: 0
- `HandleGetWeeklySummary()` - handler:
  1. Calculate week start (Monday 00:00:00) and end (Sunday 23:59:59)
  2. Fetch current user
  3. Fetch all time entries for the week
  4. Aggregate by project
  5. Optionally show daily breakdown

### 2. Helper functions
Reuse `parseDuration()` from daily summary.

Add week calculation helper:
```go
func getWeekBounds(offset int) (start, end time.Time)
```
- Returns Monday 00:00:00 and Sunday 23:59:59 for the given week offset

### 3. `main.go`
Register: `s.AddTool(clockify.GetWeeklySummaryTool(), tools.HandleGetWeeklySummary)`

## Output Format
```
Weekly Summary (Jan 13 - Jan 19, 2024):

Total: 42h 30m

By Project:
- professional-work: 32h 00m (75.3%)
- main-personal-project: 8h 30m (20.0%)
- sport: 2h 00m (4.7%)

Daily Breakdown:
- Mon: 8h 30m
- Tue: 9h 00m
- Wed: 8h 15m
- Thu: 8h 45m
- Fri: 8h 00m
- Sat: 0h 00m
- Sun: 0h 00m
```

## Verification
1. Test with current week (default)
2. Test with `week_offset: -1` for last week
3. Verify totals match Clockify weekly report
