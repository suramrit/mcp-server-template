package tools

import (
	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

type Tool interface {
	GetTool() mcp.Tool
	GetHandler() server.ToolHandlerFunc
}
