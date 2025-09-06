package tools

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// newEchoTool creates our simple but lovable echo tool
// it's like a friendly parrot that repeats what you say!
func NewEchoTool() *EchoTool {
	return &EchoTool{}
}

// echoTool is our simplest tool - it just says hello back
// every good MCP server needs at least one friendly tool
type EchoTool struct {
	// no fields needed - we're keeping it simple!
}

// getTool defines what our echo tool looks like to the outside world
// this is like writing the instruction manual for our tool
func (t *EchoTool) GetTool() mcp.Tool {
	return mcp.NewTool("echo",
		mcp.WithDescription("Echo back the provided text"), // be descriptive - AI needs to know what we do!
		mcp.WithString("name",
			mcp.Required(), // this parameter is mandatory - no anonymous greetings!
			mcp.Description("Name of the person to greet"), // help text for the AI
		),
	)
}

// getHandler returns the actual function that does the work
// this is where the rubber meets the road!
func (t *EchoTool) GetHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// extract the name parameter - this could fail if AI misbehaves
		name, err := req.RequireString("name")
		if err != nil {
			// return a nice error message instead of crashing
			return mcp.NewToolResultError(err.Error()), nil
		}

		// do our incredibly sophisticated work: say hello!
		return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
	}
}
