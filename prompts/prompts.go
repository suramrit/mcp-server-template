package prompts

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// prompt interface defines what every prompt must provide
// like a script template - every prompt needs these components
type Prompt interface {
	GetPrompt() mcp.Prompt                  // what kind of prompt are you? (name, description, arguments)
	GetHandler() server.PromptHandlerFunc   // how do we generate you? (the template creation logic)
}
