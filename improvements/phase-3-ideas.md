# Phase 3 Ideas

Future enhancements and optimizations for consideration.

---

## Production Hardening

### 1. Testing
- Unit tests for duration parsing (`parseDuration`)
- Unit tests for week bounds calculation (`getWeekBounds`)
- Integration tests with mock Clockify API
- Test edge cases: empty responses, invalid dates, running timers

### 2. Reliability
- Retry logic with exponential backoff for API calls
- Request timeouts (currently unlimited)
- Rate limiting handling (Clockify API has limits)
- Circuit breaker pattern for external API calls

### 3. Performance
- Cache current user info (fetched on every tool call)
- Cache project list (rarely changes, fetched for every summary)
- Pagination handling for users with many entries
- Concurrent fetching for multi-week summaries

---

## Feature Enhancements

### 1. Summary Improvements
- Filter summaries by tags
- Monthly/yearly summaries
- Export summaries to CSV/JSON
- Compare week-over-week trends
- Billable vs non-billable breakdowns

### 2. Entry Management
- Update/edit existing time entries
- Bulk operations (delete multiple, update tags)
- Clone/duplicate entries
- Templates for recurring tasks

### 3. Reporting
- Generate weekly timesheets
- Invoice-ready reports (billable hours)
- Productivity insights (peak hours, project patterns)
- Integration with calendar apps

---

## Code Quality

### 1. Refactoring Opportunities
- Extract common HTTP request logic from `Get`, `Post`, `Patch`, `Delete` into `doRequest` helper (~90% code duplication)
- Create `projectLookup` helper to encapsulate project name resolution
- Consider builder pattern for `TimeEntriesParams`

### 2. Logging
- Add structured logging with `log/slog`
- Log API requests (with sanitized keys)
- Debug mode for troubleshooting

### 3. Configuration
- Support `CLOCKIFY_WORKSPACE_ID` env var to override default
- Configurable pagination size
- Timezone configuration

---

## Developer Experience

### 1. CI/CD
- GitHub Actions for `go build`, `go test`, `go vet`
- Automated releases with goreleaser
- Docker image for easy deployment
- Pre-commit hooks

### 2. Documentation
- API documentation with examples
- Troubleshooting guide
- Contributing guidelines
- Architecture decision records (ADRs)

### 3. Tooling
- Makefile for common tasks
- Development setup script
- Mock Clockify API for local testing
