package clockify

import (
	"context"
	"fmt"
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