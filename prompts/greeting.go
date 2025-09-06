package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// greetingPrompt provides AI with a friendly conversation template
// think of this as giving AI a script for how to start conversations
type GreetingPrompt struct {
	// no state needed - we're just a template generator
}

// newGreetingPrompt creates our conversation-starting prompt
// every good AI assistant needs to know how to be polite!
func NewGreetingPrompt() *GreetingPrompt {
	return &GreetingPrompt{}
}

// getPrompt defines what our greeting prompt looks like
// this is like writing the description for a conversation template
func (p *GreetingPrompt) GetPrompt() mcp.Prompt {
	return mcp.NewPrompt("greeting",
		mcp.WithPromptDescription("A friendly greeting prompt"), // tell AI what this template is for
		mcp.WithArgument("name",
			mcp.ArgumentDescription("Name of the person to greet"), // optional parameter for personalization
		),
	)
}

// getHandler returns the function that generates the actual prompt
// this is where we create the conversation template on demand
func (p *GreetingPrompt) GetHandler() server.PromptHandlerFunc {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		// extract the name parameter if provided, or use a default
		name := request.Params.Arguments["name"]
		if name == "" {
			name = "friend" // everyone needs a friend!
		}

		// create a structured prompt that AI can use as a conversation starter
		return mcp.NewGetPromptResult(
			"A friendly greeting", // description of what we're returning
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant, // this message comes from the AI assistant
					mcp.NewTextContent(fmt.Sprintf("Hello, %s! How can I help you today?", name)),
				),
			},
		), nil
	}
}
