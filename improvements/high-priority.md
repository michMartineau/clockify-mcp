# High Priority Improvements

Issues that significantly affect functionality or user experience. Recommended to fix before production use.

---

## 1. Running Timers Show Zero Duration in Summaries

**Priority:** High
**Effort:** Medium
**Affected Files:** `clockify/tools.go:556-565`, `650-667`

### Current Behavior
Running timers have an empty duration string (`""`), which `parseDuration` returns as `0`. If a user has been tracking time for 3 hours, the daily/weekly summary shows 0 time for that entry.

### Impact
Silent data loss - summaries don't reflect actual tracked time for active timers.

### Proposed Fix
Calculate duration for running entries based on `start` to `now`:

```go
// Add helper function in tools.go
func getEntryDuration(entry TimeEntry) time.Duration {
    duration := parseDuration(entry.TimeInterval.Duration)
    if duration == 0 && entry.TimeInterval.End == "" {
        // Running timer - calculate duration from start
        startTime, err := time.Parse(time.RFC3339, entry.TimeInterval.Start)
        if err == nil {
            duration = time.Since(startTime)
        }
    }
    return duration
}
```

Then replace `parseDuration(entry.TimeInterval.Duration)` with `getEntryDuration(entry)` in:
- `HandleGetDailySummary` (line ~565)
- `HandleGetWeeklySummary` (line ~660)

---

## 2. Non-Deterministic Project Order in Summary Output

**Priority:** High
**Effort:** Low
**Affected Files:** `clockify/tools.go:574-586`, `677-689`

### Current Behavior
Projects appear in random order each time (Go map iteration is non-deterministic).

### Impact
Difficult to compare summaries or scan for specific projects.

### Proposed Fix
Sort projects by duration (descending) before displaying:

```go
// After aggregating durations, before output formatting
type projectSummary struct {
    ID       string
    Name     string
    Duration time.Duration
}

var summaries []projectSummary
for id, dur := range projectDurations {
    name := projectNames[id]
    if name == "" {
        name = id
    }
    summaries = append(summaries, projectSummary{ID: id, Name: name, Duration: dur})
}

// Sort by duration descending (highest first)
sort.Slice(summaries, func(i, j int) bool {
    return summaries[i].Duration > summaries[j].Duration
})

// Format output using sorted slice
for _, s := range summaries {
    percentage := 0.0
    if totalDuration > 0 {
        percentage = (float64(s.Duration) / float64(totalDuration)) * 100
    }
    result += fmt.Sprintf("- %s: %s (%.1f%%)\n", s.Name, formatDuration(s.Duration), percentage)
}
```

---

## 3. Silent Failure When Fetching Project Names

**Priority:** High
**Effort:** Low
**Affected Files:** `clockify/tools.go:549-554`, `631-638`

### Current Behavior
When `GetProjects` fails, error is silently ignored. Summaries fall back to showing project IDs.

### Impact
Masks underlying issues (rate limiting, auth problems) and degrades UX without user awareness.

### Proposed Fix
Add note to output when project name resolution fails:

```go
projects, err := t.client.GetProjects(ctx, user.DefaultWorkspace)
projectNamesFailed := false
if err != nil {
    projectNamesFailed = true
} else {
    for _, p := range projects {
        projectNames[p.ID] = p.Name
    }
}

// Later in the output formatting section:
if projectNamesFailed {
    result += "\n(Note: Could not resolve project names - showing IDs)\n"
}
```

**Alternative:** Return error and fail the summary if project fetching fails (more strict).
