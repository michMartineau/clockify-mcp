# Plan: Implement `get_daily_summary` MCP Tool

## Overview
Add the `get_daily_summary` tool to aggregate time by project for a given day.

## API Details
- **Uses:** `GET /workspaces/{workspaceId}/user/{userId}/time-entries?start=...&end=...`
- **Logic:** Fetch all entries for the day, aggregate durations by project

## Files to Modify

### 1. `clockify/client.go`
No new client methods needed - reuse `GetTimeEntries()`.

### 2. `clockify/tools.go`
- `GetDailySummaryTool()` - tool definition with parameters:
  - `date` (optional) - ISO 8601 date (YYYY-MM-DD), defaults to today
- `HandleGetDailySummary()` - handler:
  1. Parse date, compute start (00:00:00) and end (23:59:59)
  2. Fetch current user
  3. Fetch all time entries for that day
  4. Aggregate by project:
     - Parse ISO 8601 durations (PT1H30M format)
     - Sum per project
  5. Format output: project name, total time, percentage of day

### 3. Helper function
Add duration parsing helper:
```go
func parseDuration(isoDuration string) time.Duration
```
- Parses "PT1H30M45S" format to Go duration
- Handle edge cases (running entries with no duration)

### 4. `main.go`
Register: `s.AddTool(clockify.GetDailySummaryTool(), tools.HandleGetDailySummary)`

## Output Format
```
Daily Summary for 2024-01-15:

Total: 8h 30m

By Project:
- professional-work: 6h 15m (73.5%)
- main-personal-project: 2h 15m (26.5%)
```

## Verification
1. Test with today's date (default)
2. Test with explicit past date
3. Verify totals match Clockify dashboard
