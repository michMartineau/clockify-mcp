package main

import (
	"log"

	"clockify-mcp/clockify"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	client, err := clockify.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Clockify client: %v", err)
	}

	tools := clockify.NewTools(client)

	s := server.NewMCPServer(
		"clockify-mcp",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(clockify.GetCurrentUserTool(), tools.HandleGetCurrentUser)
	s.AddTool(clockify.ListProjectsTool(), tools.HandleListProjects)
	s.AddTool(clockify.ListTagsTool(), tools.HandleListTags)
	s.AddTool(clockify.ListTimeEntriesTool(), tools.HandleListTimeEntries)
	s.AddTool(clockify.GetCurrentTimerTool(), tools.HandleGetCurrentTimer)
	s.AddTool(clockify.StartTimerTool(), tools.HandleStartTimer)
	s.AddTool(clockify.StopTimerTool(), tools.HandleStopTimer)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}