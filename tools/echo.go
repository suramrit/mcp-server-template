package tools

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

func NewEchoTool() *EchoTool {
	return &EchoTool{}
}

type EchoTool struct {
}

func (t *EchoTool) GetTool() mcp.Tool {
	return mcp.NewTool("echo",
		mcp.WithDescription("Echo back the provided text"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
}

func (t *EchoTool) GetHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
	}
}
