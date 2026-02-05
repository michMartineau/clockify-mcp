# Plan: Implement `list_projects` MCP Tool

## Overview
Add the `list_projects` tool to the Clockify MCP server, following the existing pattern established by `get_current_user`.

## API Details
- **Endpoint:** `GET /workspaces/{workspaceId}/projects`
- **Parameters:**
  - `workspaceId` (required) - defaults to user's default workspace if not provided
  - `page` (optional) - pagination, default 1
  - `page-size` (optional) - default 50

## Files to Modify

### 1. `clockify/models.go`
Add `Project` struct matching Clockify API response:
```go
type Project struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    ClientID   string `json:"clientId"`
    ClientName string `json:"clientName"`
    Archived   bool   `json:"archived"`
    Billable   bool   `json:"billable"`
    Color      string `json:"color"`
}
```

### 2. `clockify/client.go`
Add method to fetch projects:
```go
func (c *Client) GetProjects(ctx context.Context, workspaceID string) ([]Project, error)
```
- Calls `GET /workspaces/{workspaceID}/projects`
- Returns slice of projects

### 3. `clockify/tools.go`
Add tool definition and handler:
- `ListProjectsTool()` - defines the MCP tool with optional `workspace_id` parameter
- `HandleListProjects()` - handles the tool call:
  1. Extract optional `workspace_id` from request
  2. If not provided, fetch current user to get default workspace
  3. Call `client.GetProjects()`
  4. Format and return results

### 4. `main.go`
Register the new tool:
```go
s.AddTool(clockify.ListProjectsTool(), tools.HandleListProjects)
```

## Implementation Order
1. Add `Project` model to `models.go`
2. Add `GetProjects` client method to `client.go`
3. Add tool definition and handler to `tools.go`
4. Register tool in `main.go`

## Verification
1. Build: `go build -o clockify-mcp .`
2. Test manually with Claude Code MCP integration or stdio
3. Verify tool lists projects from default workspace when no `workspace_id` provided
4. Verify tool accepts optional `workspace_id` parameter