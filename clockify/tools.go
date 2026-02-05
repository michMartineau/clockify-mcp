package clockify

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// Tools holds the MCP tool handlers for Clockify.
type Tools struct {
	client *Client
}

// NewTools creates a new Tools instance.
func NewTools(client *Client) *Tools {
	return &Tools{client: client}
}

// GetCurrentUserTool returns the MCP tool definition for get_current_user.
func GetCurrentUserTool() mcp.Tool {
	return mcp.NewTool(
		"get_current_user",
		mcp.WithDescription("Get the current Clockify user's information including user ID, email, name, and default workspace ID"),
	)
}

// HandleGetCurrentUser handles the get_current_user tool call.
func (t *Tools) HandleGetCurrentUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	result := fmt.Sprintf("User ID: %s\nEmail: %s\nName: %s\nDefault Workspace ID: %s",
		user.ID, user.Email, user.Name, user.DefaultWorkspace)

	return mcp.NewToolResultText(result), nil
}

// ListProjectsTool returns the MCP tool definition for list_projects.
func ListProjectsTool() mcp.Tool {
	return mcp.NewTool(
		"list_projects",
		mcp.WithDescription("List all projects in a Clockify workspace"),
		mcp.WithString("workspace_id",
			mcp.Description("The workspace ID. If not provided, uses the current user's default workspace"),
		),
	)
}

// HandleListProjects handles the list_projects tool call.
func (t *Tools) HandleListProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workspaceID := mcp.ParseString(request, "workspace_id", "")

	// If no workspace ID provided, get the default from current user
	if workspaceID == "" {
		user, err := t.client.GetCurrentUser(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
		}
		workspaceID = user.DefaultWorkspace
	}

	projects, err := t.client.GetProjects(ctx, workspaceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get projects: %v", err)), nil
	}

	if len(projects) == 0 {
		return mcp.NewToolResultText("No projects found in workspace"), nil
	}

	result := fmt.Sprintf("Found %d projects:\n\n", len(projects))
	for _, p := range projects {
		archived := ""
		if p.Archived {
			archived = " [ARCHIVED]"
		}
		result += fmt.Sprintf("- %s (ID: %s)%s\n", p.Name, p.ID, archived)
		if p.ClientName != "" {
			result += fmt.Sprintf("  Client: %s\n", p.ClientName)
		}
	}

	return mcp.NewToolResultText(result), nil
}

// ListTagsTool returns the MCP tool definition for list_tags.
func ListTagsTool() mcp.Tool {
	return mcp.NewTool(
		"list_tags",
		mcp.WithDescription("List all tags in a Clockify workspace"),
		mcp.WithString("workspace_id",
			mcp.Description("The workspace ID. If not provided, uses the current user's default workspace"),
		),
	)
}

// HandleListTags handles the list_tags tool call.
func (t *Tools) HandleListTags(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workspaceID := mcp.ParseString(request, "workspace_id", "")

	// If no workspace ID provided, get the default from current user
	if workspaceID == "" {
		user, err := t.client.GetCurrentUser(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
		}
		workspaceID = user.DefaultWorkspace
	}

	tags, err := t.client.GetTags(ctx, workspaceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tags: %v", err)), nil
	}

	if len(tags) == 0 {
		return mcp.NewToolResultText("No tags found in workspace"), nil
	}

	result := fmt.Sprintf("Found %d tags:\n\n", len(tags))
	for _, tag := range tags {
		archived := ""
		if tag.Archived {
			archived = " [ARCHIVED]"
		}
		result += fmt.Sprintf("- %s (ID: %s)%s\n", tag.Name, tag.ID, archived)
	}

	return mcp.NewToolResultText(result), nil
}

// ListTimeEntriesTool returns the MCP tool definition for list_time_entries.
func ListTimeEntriesTool() mcp.Tool {
	return mcp.NewTool(
		"list_time_entries",
		mcp.WithDescription("List time entries for a date range. Returns time entries with description, project, start/end times, and duration."),
		mcp.WithString("start_date",
			mcp.Description("Start date/time in ISO 8601 format (e.g., 2024-01-15T00:00:00Z). Defaults to start of today (00:00:00)"),
		),
		mcp.WithString("end_date",
			mcp.Description("End date/time in ISO 8601 format (e.g., 2024-01-15T23:59:59Z). Defaults to now"),
		),
		mcp.WithString("project_id",
			mcp.Description("Optional project ID to filter entries by"),
		),
	)
}

// HandleListTimeEntries handles the list_time_entries tool call.
func (t *Tools) HandleListTimeEntries(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startDate := mcp.ParseString(request, "start_date", "")
	endDate := mcp.ParseString(request, "end_date", "")
	projectID := mcp.ParseString(request, "project_id", "")

	// Get current user for workspace and user IDs
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	// Default start_date to start of today (00:00:00 UTC)
	if startDate == "" {
		now := time.Now().UTC()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		startDate = startOfDay.Format(time.RFC3339)
	}

	// Default end_date to now
	if endDate == "" {
		endDate = time.Now().UTC().Format(time.RFC3339)
	}

	params := TimeEntriesParams{
		Start:     startDate,
		End:       endDate,
		ProjectID: projectID,
	}

	entries, err := t.client.GetTimeEntries(ctx, user.DefaultWorkspace, user.ID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get time entries: %v", err)), nil
	}

	if len(entries) == 0 {
		return mcp.NewToolResultText("No time entries found for the specified date range"), nil
	}

	result := fmt.Sprintf("Found %d time entries:\n\n", len(entries))
	for _, entry := range entries {
		description := entry.Description
		if description == "" {
			description = "(no description)"
		}

		// Format duration for display
		duration := entry.TimeInterval.Duration
		if duration == "" {
			duration = "running"
		}

		result += fmt.Sprintf("- %s\n", description)
		result += fmt.Sprintf("  ID: %s\n", entry.ID)
		result += fmt.Sprintf("  Start: %s\n", entry.TimeInterval.Start)
		if entry.TimeInterval.End != "" {
			result += fmt.Sprintf("  End: %s\n", entry.TimeInterval.End)
		}
		result += fmt.Sprintf("  Duration: %s\n", duration)
		if entry.ProjectID != "" {
			result += fmt.Sprintf("  Project ID: %s\n", entry.ProjectID)
		}
		result += "\n"
	}

	return mcp.NewToolResultText(result), nil
}

// GetCurrentTimerTool returns the MCP tool definition for get_current_timer.
func GetCurrentTimerTool() mcp.Tool {
	return mcp.NewTool(
		"get_current_timer",
		mcp.WithDescription("Get the currently running Clockify time entry, if any"),
	)
}

// HandleGetCurrentTimer handles the get_current_timer tool call.
func (t *Tools) HandleGetCurrentTimer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get current user to get workspace and user IDs
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	entry, err := t.client.GetCurrentTimer(ctx, user.DefaultWorkspace, user.ID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current timer: %v", err)), nil
	}

	if entry == nil {
		return mcp.NewToolResultText("No timer running"), nil
	}

	result := fmt.Sprintf("Timer running:\n")
	result += fmt.Sprintf("ID: %s\n", entry.ID)
	if entry.Description != "" {
		result += fmt.Sprintf("Description: %s\n", entry.Description)
	}
	result += fmt.Sprintf("Started: %s\n", entry.TimeInterval.Start)
	if entry.ProjectID != "" {
		result += fmt.Sprintf("Project ID: %s\n", entry.ProjectID)
	}
	if len(entry.TagIDs) > 0 {
		result += fmt.Sprintf("Tag IDs: %v\n", entry.TagIDs)
	}
	result += fmt.Sprintf("Billable: %t\n", entry.Billable)

	return mcp.NewToolResultText(result), nil
}

// StartTimerTool returns the MCP tool definition for start_timer.
func StartTimerTool() mcp.Tool {
	return mcp.NewTool(
		"start_timer",
		mcp.WithDescription("Start a new Clockify time entry (running timer)"),
		mcp.WithString("description",
			mcp.Description("Description of what you're working on"),
		),
		mcp.WithString("project_id",
			mcp.Description("Project ID to associate with the entry"),
		),
		mcp.WithBoolean("billable",
			mcp.Description("Whether the entry is billable"),
		),
	)
}

// HandleStartTimer handles the start_timer tool call.
func (t *Tools) HandleStartTimer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	description := mcp.ParseString(request, "description", "")
	projectID := mcp.ParseString(request, "project_id", "")
	billable := mcp.ParseBoolean(request, "billable", false)

	// Get current user to get default workspace
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	// Build the request
	req := StartTimerRequest{
		Start:       time.Now().UTC().Format(time.RFC3339),
		Description: description,
		ProjectID:   projectID,
		Billable:    billable,
	}

	entry, err := t.client.StartTimer(ctx, user.DefaultWorkspace, req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start timer: %v", err)), nil
	}

	result := fmt.Sprintf("Timer started:\n")
	result += fmt.Sprintf("ID: %s\n", entry.ID)
	if entry.Description != "" {
		result += fmt.Sprintf("Description: %s\n", entry.Description)
	}
	result += fmt.Sprintf("Started: %s\n", entry.TimeInterval.Start)
	if entry.ProjectID != "" {
		result += fmt.Sprintf("Project ID: %s\n", entry.ProjectID)
	}
	result += fmt.Sprintf("Billable: %t\n", entry.Billable)

	return mcp.NewToolResultText(result), nil
}

// StopTimerTool returns the MCP tool definition for stop_timer.
func StopTimerTool() mcp.Tool {
	return mcp.NewTool(
		"stop_timer",
		mcp.WithDescription("Stop the currently running Clockify time entry"),
	)
}

// HandleStopTimer handles the stop_timer tool call.
func (t *Tools) HandleStopTimer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get current user to get workspace and user IDs
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	// Stop the timer with current time
	endTime := time.Now().UTC().Format(time.RFC3339)
	entry, err := t.client.StopTimer(ctx, user.DefaultWorkspace, user.ID, endTime)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop timer: %v", err)), nil
	}

	result := fmt.Sprintf("Timer stopped:\n")
	result += fmt.Sprintf("ID: %s\n", entry.ID)
	if entry.Description != "" {
		result += fmt.Sprintf("Description: %s\n", entry.Description)
	}
	result += fmt.Sprintf("Started: %s\n", entry.TimeInterval.Start)
	result += fmt.Sprintf("Ended: %s\n", entry.TimeInterval.End)
	if entry.TimeInterval.Duration != "" {
		result += fmt.Sprintf("Duration: %s\n", entry.TimeInterval.Duration)
	}
	if entry.ProjectID != "" {
		result += fmt.Sprintf("Project ID: %s\n", entry.ProjectID)
	}
	result += fmt.Sprintf("Billable: %t\n", entry.Billable)

	return mcp.NewToolResultText(result), nil
}

// CreateTimeEntryTool returns the MCP tool definition for create_time_entry.
func CreateTimeEntryTool() mcp.Tool {
	return mcp.NewTool(
		"create_time_entry",
		mcp.WithDescription("Create a completed time entry with both start and end times. Use this to log past work."),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Description("Start time in ISO 8601 format (e.g., 2024-01-15T09:00:00Z)"),
		),
		mcp.WithString("end",
			mcp.Required(),
			mcp.Description("End time in ISO 8601 format (e.g., 2024-01-15T17:00:00Z)"),
		),
		mcp.WithString("description",
			mcp.Description("Description of what was worked on"),
		),
		mcp.WithString("project_id",
			mcp.Description("Project ID to associate with the entry"),
		),
		mcp.WithBoolean("billable",
			mcp.Description("Whether the entry is billable"),
		),
	)
}

// HandleCreateTimeEntry handles the create_time_entry tool call.
func (t *Tools) HandleCreateTimeEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := mcp.ParseString(request, "start", "")
	end := mcp.ParseString(request, "end", "")
	description := mcp.ParseString(request, "description", "")
	projectID := mcp.ParseString(request, "project_id", "")
	billable := mcp.ParseBoolean(request, "billable", false)

	// Validate required parameters
	if start == "" {
		return mcp.NewToolResultError("start parameter is required"), nil
	}
	if end == "" {
		return mcp.NewToolResultError("end parameter is required"), nil
	}

	// Validate start < end
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid start time format: %v", err)), nil
	}
	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid end time format: %v", err)), nil
	}
	if !startTime.Before(endTime) {
		return mcp.NewToolResultError("Start time must be before end time"), nil
	}

	// Get current user to get default workspace
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	// Build the request
	req := CreateTimeEntryRequest{
		Start:       start,
		End:         end,
		Description: description,
		ProjectID:   projectID,
		Billable:    billable,
	}

	entry, err := t.client.CreateTimeEntry(ctx, user.DefaultWorkspace, req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create time entry: %v", err)), nil
	}

	result := fmt.Sprintf("Time entry created:\n")
	result += fmt.Sprintf("ID: %s\n", entry.ID)
	if entry.Description != "" {
		result += fmt.Sprintf("Description: %s\n", entry.Description)
	}
	result += fmt.Sprintf("Started: %s\n", entry.TimeInterval.Start)
	result += fmt.Sprintf("Ended: %s\n", entry.TimeInterval.End)
	if entry.TimeInterval.Duration != "" {
		result += fmt.Sprintf("Duration: %s\n", entry.TimeInterval.Duration)
	}
	if entry.ProjectID != "" {
		result += fmt.Sprintf("Project ID: %s\n", entry.ProjectID)
	}
	result += fmt.Sprintf("Billable: %t\n", entry.Billable)

	return mcp.NewToolResultText(result), nil
}

// DeleteTimeEntryTool returns the MCP tool definition for delete_time_entry.
func DeleteTimeEntryTool() mcp.Tool {
	return mcp.NewTool(
		"delete_time_entry",
		mcp.WithDescription("Delete a time entry by ID"),
		mcp.WithString("entry_id",
			mcp.Required(),
			mcp.Description("The ID of the time entry to delete"),
		),
	)
}

// HandleDeleteTimeEntry handles the delete_time_entry tool call.
func (t *Tools) HandleDeleteTimeEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	entryID := mcp.ParseString(request, "entry_id", "")

	// Validate required parameter
	if entryID == "" {
		return mcp.NewToolResultError("entry_id parameter is required"), nil
	}

	// Get current user to get default workspace
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	// Delete the time entry
	err = t.client.DeleteTimeEntry(ctx, user.DefaultWorkspace, entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete time entry: %v", err)), nil
	}

	result := fmt.Sprintf("Time entry deleted successfully (ID: %s)", entryID)
	return mcp.NewToolResultText(result), nil
}

// GetDailySummaryTool returns the MCP tool definition for get_daily_summary.
func GetDailySummaryTool() mcp.Tool {
	return mcp.NewTool(
		"get_daily_summary",
		mcp.WithDescription("Get a daily summary of time tracked by project for a specific date. Returns total time and breakdown by project with percentages."),
		mcp.WithString("date",
			mcp.Description("Date in ISO 8601 format (YYYY-MM-DD, e.g., 2024-01-15). Defaults to today"),
		),
	)
}

// HandleGetDailySummary handles the get_daily_summary tool call.
func (t *Tools) HandleGetDailySummary(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dateStr := mcp.ParseString(request, "date", "")

	// Parse date or default to today
	var targetDate time.Time
	if dateStr == "" {
		targetDate = time.Now().UTC()
	} else {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid date format. Use YYYY-MM-DD: %v", err)), nil
		}
		targetDate = parsed.UTC()
	}

	// Calculate start (00:00:00) and end (23:59:59) of the day
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 999999999, time.UTC)

	// Get current user
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	// Fetch time entries for the day
	params := TimeEntriesParams{
		Start: startOfDay.Format(time.RFC3339),
		End:   endOfDay.Format(time.RFC3339),
	}

	entries, err := t.client.GetTimeEntries(ctx, user.DefaultWorkspace, user.ID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get time entries: %v", err)), nil
	}

	if len(entries) == 0 {
		dateDisplay := targetDate.Format("2006-01-02")
		return mcp.NewToolResultText(fmt.Sprintf("Daily Summary for %s:\n\nNo time entries found", dateDisplay)), nil
	}

	// Aggregate by project
	projectDurations := make(map[string]time.Duration)
	projectNames := make(map[string]string) // projectID -> projectName
	var totalDuration time.Duration

	// Get all projects to map IDs to names
	projects, err := t.client.GetProjects(ctx, user.DefaultWorkspace)
	if err == nil {
		for _, p := range projects {
			projectNames[p.ID] = p.Name
		}
	}

	for _, entry := range entries {
		duration := parseDuration(entry.TimeInterval.Duration)
		totalDuration += duration

		projectKey := entry.ProjectID
		if projectKey == "" {
			projectKey = "(no project)"
		}
		projectDurations[projectKey] += duration
	}

	// Format output
	dateDisplay := targetDate.Format("2006-01-02")
	result := fmt.Sprintf("Daily Summary for %s:\n\n", dateDisplay)
	result += fmt.Sprintf("Total: %s\n\n", formatDuration(totalDuration))

	if len(projectDurations) > 0 {
		result += "By Project:\n"
		for projectID, duration := range projectDurations {
			projectName := projectNames[projectID]
			if projectName == "" {
				projectName = projectID
			}

			percentage := 0.0
			if totalDuration > 0 {
				percentage = (float64(duration) / float64(totalDuration)) * 100
			}

			result += fmt.Sprintf("- %s: %s (%.1f%%)\n", projectName, formatDuration(duration), percentage)
		}
	}

	return mcp.NewToolResultText(result), nil
}

// GetWeeklySummaryTool returns the MCP tool definition for get_weekly_summary.
func GetWeeklySummaryTool() mcp.Tool {
	return mcp.NewTool(
		"get_weekly_summary",
		mcp.WithDescription("Get a weekly summary of time tracked by project. Returns total time, breakdown by project with percentages, and daily breakdown. Week runs Monday-Sunday."),
		mcp.WithString("week_offset",
			mcp.Description("Week offset from current week. 0 = current week, -1 = last week, -2 = two weeks ago, etc. Defaults to 0"),
		),
	)
}

// HandleGetWeeklySummary handles the get_weekly_summary tool call.
func (t *Tools) HandleGetWeeklySummary(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	weekOffsetStr := mcp.ParseString(request, "week_offset", "0")
	weekOffset, err := strconv.Atoi(weekOffsetStr)
	if err != nil {
		weekOffset = 0
	}

	// Calculate week bounds (Monday 00:00:00 to Sunday 23:59:59)
	weekStart, weekEnd := getWeekBounds(weekOffset)

	// Get current user
	user, err := t.client.GetCurrentUser(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}

	// Fetch time entries for the week
	params := TimeEntriesParams{
		Start: weekStart.Format(time.RFC3339),
		End:   weekEnd.Format(time.RFC3339),
	}

	entries, err := t.client.GetTimeEntries(ctx, user.DefaultWorkspace, user.ID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get time entries: %v", err)), nil
	}

	// Get all projects to map IDs to names
	projects, err := t.client.GetProjects(ctx, user.DefaultWorkspace)
	projectNames := make(map[string]string)
	if err == nil {
		for _, p := range projects {
			projectNames[p.ID] = p.Name
		}
	}

	if len(entries) == 0 {
		weekDisplay := fmt.Sprintf("%s - %s", weekStart.Format("Jan 02"), weekEnd.Format("Jan 02, 2006"))
		return mcp.NewToolResultText(fmt.Sprintf("Weekly Summary (%s):\n\nNo time entries found", weekDisplay)), nil
	}

	// Aggregate by project and by day
	projectDurations := make(map[string]time.Duration)
	dailyDurations := make(map[string]time.Duration) // day name -> duration
	var totalDuration time.Duration

	for _, entry := range entries {
		duration := parseDuration(entry.TimeInterval.Duration)
		totalDuration += duration

		// Aggregate by project
		projectKey := entry.ProjectID
		if projectKey == "" {
			projectKey = "(no project)"
		}
		projectDurations[projectKey] += duration

		// Aggregate by day
		startTime, err := time.Parse(time.RFC3339, entry.TimeInterval.Start)
		if err == nil {
			dayName := startTime.Weekday().String()[:3] // Mon, Tue, etc.
			dailyDurations[dayName] += duration
		}
	}

	// Format output
	weekDisplay := fmt.Sprintf("%s - %s", weekStart.Format("Jan 02"), weekEnd.Format("Jan 02, 2006"))
	result := fmt.Sprintf("Weekly Summary (%s):\n\n", weekDisplay)
	result += fmt.Sprintf("Total: %s\n\n", formatDuration(totalDuration))

	// By project
	if len(projectDurations) > 0 {
		result += "By Project:\n"
		for projectID, duration := range projectDurations {
			projectName := projectNames[projectID]
			if projectName == "" {
				projectName = projectID
			}

			percentage := 0.0
			if totalDuration > 0 {
				percentage = (float64(duration) / float64(totalDuration)) * 100
			}

			result += fmt.Sprintf("- %s: %s (%.1f%%)\n", projectName, formatDuration(duration), percentage)
		}
	}

	// Daily breakdown
	result += "\nDaily Breakdown:\n"
	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	for _, day := range weekdays {
		duration := dailyDurations[day]
		result += fmt.Sprintf("- %s: %s\n", day, formatDuration(duration))
	}

	return mcp.NewToolResultText(result), nil
}

// parseDuration parses an ISO 8601 duration string (e.g., "PT1H30M45S") into time.Duration.
// Returns 0 for empty or invalid durations (e.g., running timers).
func parseDuration(isoDuration string) time.Duration {
	if isoDuration == "" {
		return 0
	}

	// Parse ISO 8601 duration format: PT1H30M45S
	// Pattern: PT([0-9]+H)?([0-9]+M)?([0-9]+S)?
	re := regexp.MustCompile(`^PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?$`)
	matches := re.FindStringSubmatch(isoDuration)

	if matches == nil {
		return 0
	}

	var duration time.Duration

	// Hours
	if matches[1] != "" {
		hours, _ := strconv.Atoi(matches[1])
		duration += time.Duration(hours) * time.Hour
	}

	// Minutes
	if matches[2] != "" {
		minutes, _ := strconv.Atoi(matches[2])
		duration += time.Duration(minutes) * time.Minute
	}

	// Seconds (can be decimal)
	if matches[3] != "" {
		seconds, _ := strconv.ParseFloat(matches[3], 64)
		duration += time.Duration(seconds * float64(time.Second))
	}

	return duration
}

// getWeekBounds calculates the start and end times for a week given an offset.
// offset=0 is current week, offset=-1 is last week, etc.
// Week starts on Monday 00:00:00 and ends on Sunday 23:59:59.
func getWeekBounds(offset int) (start, end time.Time) {
	now := time.Now().UTC()

	// Find Monday of the current week
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	daysFromMonday := weekday - 1

	// Calculate Monday of target week
	monday := now.AddDate(0, 0, -daysFromMonday+(offset*7))
	start = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate Sunday of target week
	sunday := start.AddDate(0, 0, 6)
	end = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, time.UTC)

	return start, end
}

// formatDuration formats a duration as "Xh Ym" (e.g., "2h 30m").
// If duration is 0, returns "0h 0m".
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0h 0m"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	return fmt.Sprintf("%dh %dm", hours, minutes)
}