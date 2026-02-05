# Plan: Implement `list_tags` MCP Tool

## Overview
Add the `list_tags` tool to list available tags in a workspace.

## API Details
- **Endpoint:** `GET /workspaces/{workspaceId}/tags`
- **Parameters:**
  - `workspaceId` (required) - defaults to user's default workspace
  - `page`, `page-size` (optional) - pagination

## Files to Modify

### 1. `clockify/models.go`
Add `Tag` struct:
```go
type Tag struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    WorkspaceID string `json:"workspaceId"`
    Archived    bool   `json:"archived"`
}
```

### 2. `clockify/client.go`
Add method:
```go
func (c *Client) GetTags(ctx context.Context, workspaceID string) ([]Tag, error)
```

### 3. `clockify/tools.go`
- `ListTagsTool()` - tool definition with optional `workspace_id` parameter
- `HandleListTags()` - handler that fetches tags, uses default workspace if not provided

### 4. `main.go`
Register: `s.AddTool(clockify.ListTagsTool(), tools.HandleListTags)`

## Verification
1. Build: `go build -o clockify-mcp .`
2. Test tool returns tags from default workspace
3. Test with explicit `workspace_id` parameter