# Clockify MCP Server

A Go-based [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server that integrates with the [Clockify](https://clockify.me) time tracking API. Exposes 11 MCP tools for managing time entries, projects, tags, timers, and summaries over stdio transport.

## Requirements

- Go 1.21+
- Clockify account and API key

## Installation

```bash
git clone https://github.com/your-username/clockify-mcp
cd clockify-mcp
make build
```

The binary is output to `./clockify-mcp`.

## Configuration

Set your Clockify API key (found in **Profile Settings → API** in Clockify):

```bash
export CLOCKIFY_API_KEY="your-api-key-here"
```

## MCP Server Setup

### Claude Desktop

Add the following to your `claude_desktop_config.json` (usually at `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "clockify": {
      "command": "/path/to/clockify-mcp",
      "env": {
        "CLOCKIFY_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

Restart Claude Desktop after saving the file.

### Claude Code

Run the following command to register the server:

```bash
claude mcp add clockify /path/to/clockify-mcp -e CLOCKIFY_API_KEY=your-api-key-here
```

Or add it to your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "clockify": {
      "command": "/path/to/clockify-mcp",
      "env": {
        "CLOCKIFY_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

Verify the server is registered with:

```bash
claude mcp list
```

## Tools

| Tool | Description |
|------|-------------|
| `get_current_user` | Get current user info (ID, email, name, workspace) |
| `list_projects` | List all projects in a workspace |
| `list_tags` | List all tags in a workspace |
| `list_time_entries` | List time entries for a date range |
| `get_current_timer` | Get the currently running timer, if any |
| `start_timer` | Start a new running time entry |
| `stop_timer` | Stop the currently running timer |
| `create_time_entry` | Create a completed time entry with start and end times |
| `delete_time_entry` | Delete a time entry by ID |
| `get_daily_summary` | Daily breakdown of time tracked by project |
| `get_weekly_summary` | Weekly breakdown with per-project and per-day totals |

All tools auto-resolve the workspace ID from your user profile when not explicitly provided. Timestamps use ISO 8601 format and are interpreted as UTC.

## Development

```bash
make build    # Build binary
make test     # Run tests
make vet      # Static analysis
make fmt      # Format code
make check    # Run fmt → vet → test → build
make clean    # Remove binary
make deps     # Download and tidy dependencies
```

## Architecture

The server is structured in three layers under the `clockify/` package:

- **`client.go`** — HTTP client with `X-Api-Key` auth, wrapping `net/http`. Base URL: `https://api.clockify.me/api/v1`
- **`models.go`** — Data structures (`User`, `Project`, `Tag`, `TimeEntry`, etc.) and request DTOs
- **`tools.go`** — MCP tool definitions (`*Tool()` methods) and handlers (`Handle*()` methods)

`main.go` wires everything together: reads `CLOCKIFY_API_KEY`, registers all 11 tools, and serves via stdio.