package prompts

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Prompt interface {
	GetPrompt() mcp.Prompt
	GetHandler() server.PromptHandlerFunc
}
