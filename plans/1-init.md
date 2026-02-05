Set up a new MCP server project for Clockify time tracking.

1. Initialize the project:
    - Create the directory structure as defined in CLAUDE.md
    - Initialize go modules
    - Add mcp-go dependency: github.com/mark3labs/mcp-go

2. Implement the Clockify HTTP client (clockify/client.go):
    - Read API key from CLOCKIFY_API_KEY env var
    - Base URL: https://api.clockify.me/api/v1
    - Set X-Api-Key header on all requests
    - Helper methods for GET, POST, PATCH with JSON handling
    - Return errors clearly, don't swallow them

3. Implement the first tool: get_current_user
    - Calls GET /user
    - Returns user ID, email, name, and default workspace ID
    - Register it with the MCP server

4. Set up main.go:
    - Create MCP server with stdio transport
    - Register the get_current_user tool
    - Run the server

5. Test:
    - Build the binary
    - Confirm it starts without error
    - Confirm the tool is listed when inspected

Keep code simple and idiomatic Go. No external dependencies beyond mcp-go.