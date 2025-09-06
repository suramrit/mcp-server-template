package tools

import (
	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// tool interface defines what every tool must provide
// this is like a job description - if you want to be a tool, you need these skills
type Tool interface {
	GetTool() mcp.Tool                    // tell us what you are (name, description, parameters)
	GetHandler() server.ToolHandlerFunc   // show us what you can do (the actual implementation)
}
