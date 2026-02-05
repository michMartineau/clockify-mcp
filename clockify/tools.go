package clockify

import (
	"context"
	"fmt"

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