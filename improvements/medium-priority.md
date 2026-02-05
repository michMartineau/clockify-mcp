# Medium Priority Improvements

Issues that affect code quality or edge cases but don't block core functionality.

---

## 1. Invalid `week_offset` Parameter Silently Defaults to 0

**Priority:** Medium
**Effort:** Low
**Affected Files:** `clockify/tools.go:605-609`

### Current Behavior
Invalid `week_offset` like `"abc"` silently defaults to `0` (current week).

### Impact
Inconsistent with `HandleGetDailySummary` which validates date format. User receives unexpected results without error.

### Proposed Fix - Option A (Validation)
```go
weekOffsetStr := mcp.ParseString(request, "week_offset", "0")
weekOffset, err := strconv.Atoi(weekOffsetStr)
if err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("Invalid week_offset: must be an integer, got '%s'", weekOffsetStr)), nil
}
```

### Proposed Fix - Option B (Type Safety)
Change tool definition to use number type instead of string:
```go
mcp.WithNumber("week_offset",
    mcp.Description("Week offset: 0 for current week, -1 for last week, etc."),
)
```

Then parse as:
```go
weekOffset := int(mcp.ParseNumber(request, "week_offset", 0))
```

---

## 2. Unnecessary Content-Type Header on DELETE Requests

**Priority:** Medium
**Effort:** Low
**Affected Files:** `clockify/client.go:238`

### Current Behavior
`Delete` method sets `Content-Type: application/json` but sends no request body.

### Impact
Semantically incorrect (harmless but unnecessary).

### Proposed Fix
Remove Content-Type header for bodyless DELETE:

```go
func (c *Client) Delete(ctx context.Context, path string) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL+path, nil)
    if err != nil {
        return fmt.Errorf("creating request: %w", err)
    }
    req.Header.Set("X-Api-Key", c.apiKey)
    // No Content-Type needed for bodyless DELETE

    resp, err := c.httpClient.Do(req)
    // ... rest of method
}
```

---

## 3. Timezone Handling Could Surprise Users

**Priority:** Medium
**Effort:** Low (documentation)
**Affected Files:** `clockify/tools.go:497-502`

### Current Behavior
Date `"2024-01-15"` is interpreted as midnight UTC. Users in other timezones may get unexpected results.

### Impact
User in EST asking for "today" might see different date than expected.

### Proposed Fix
Document UTC behavior in tool description:

```go
mcp.WithString("date",
    mcp.Description("Date in YYYY-MM-DD format (e.g., 2024-01-15). Interpreted as midnight UTC. Defaults to today (UTC)"),
),
```

**Alternative:** Accept timezone parameter or use system timezone (more complex).

---

## 4. `Billable` Field May Not Serialize `false`

**Priority:** Medium
**Effort:** Low (needs API verification)
**Affected Files:** `clockify/models.go:63`

### Current Behavior
`Billable` field uses `omitempty`, so `false` is not included in JSON payload.

### Impact
If Clockify API defaults to `true` when field is absent, users cannot explicitly create non-billable entries.

### Investigation Needed
Test Clockify API behavior:
1. What is the default when `billable` field is omitted?
2. Does the API accept explicit `billable: false`?

### Proposed Fix (if needed)
Remove `omitempty` to always include the field:
```go
Billable    bool     `json:"billable"`
```

**OR** use pointer for explicit null vs false distinction:
```go
Billable    *bool    `json:"billable,omitempty"`
```

**Note:** This follows existing pattern in `StartTimerRequest`, so may be intentional.

---

## 5. Plan File Typo

**Priority:** Medium
**Effort:** Trivial
**Affected Files:** `plans/10-get-daily-summary.md:4`

### Current Behavior
Backticks around `get_daily_summary` were accidentally removed:
```markdown
Add the  tool to aggregate time by project for a given day.
```

### Proposed Fix
Restore original text:
```markdown
Add the `get_daily_summary` tool to aggregate time by project for a given day.
```
