package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// withToolLogging wraps tool handlers with comprehensive logging
// this is like having a security camera that records everything your tools do
func WithToolLogging(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()
		// log the incoming request - what tool is being called and with what arguments
		log.Printf("[TOOL] Starting: %s with args: %v", name, req.Params.Arguments)

		// call the actual tool handler - this is where the real work happens
		result, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			// something went wrong - log it and return a safe error message
			log.Printf("[TOOL] Error in %s (took %v): %v", name, duration, err)
			// return a sanitized error message to the client - don't leak internal details!
			return mcp.NewToolResultError(fmt.Sprintf("Tool %s failed: %v", name, err)), nil
		}

		// success! log how long it took
		log.Printf("[TOOL] Completed: %s (took %v)", name, duration)
		return result, nil
	}
}

// withResourceLogging wraps resource handlers with logging
// like a librarian who keeps track of every book that gets checked out
func WithResourceLogging(uri string, handler server.ResourceHandlerFunc) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		start := time.Now()
		log.Printf("[RESOURCE] Reading: %s", uri)

		// attempt to read the resource
		contents, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			// resource couldn't be read - log the failure
			log.Printf("[RESOURCE] Error reading %s (took %v): %v", uri, duration, err)
			return nil, fmt.Errorf("failed to read resource %s: %w", uri, err)
		}

		// success! log what we accomplished
		log.Printf("[RESOURCE] Successfully read: %s (took %v, %d items)", uri, duration, len(contents))
		return contents, nil
	}
}

// withPromptLogging wraps prompt handlers with logging
// like a director who keeps notes on every scene that gets filmed
func WithPromptLogging(name string, handler server.PromptHandlerFunc) server.PromptHandlerFunc {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		start := time.Now()
		log.Printf("[PROMPT] Getting: %s with args: %v", name, req.Params.Arguments)

		// generate the prompt template
		result, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			// prompt generation failed - this shouldn't happen often
			log.Printf("[PROMPT] Error in %s (took %v): %v", name, duration, err)
			return nil, fmt.Errorf("failed to get prompt %s: %w", name, err)
		}

		log.Printf("[PROMPT] Completed: %s (took %v)", name, duration)
		return result, nil
	}
}

// withToolRecovery catches panics in tool handlers
// this is like having a safety net under a trapeze - if something goes horribly wrong, we catch it
func WithToolRecovery(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
		// set up a panic recovery mechanism
		defer func() {
			if r := recover(); r != nil {
				// something panicked! log it and return a safe error
				log.Printf("[TOOL] PANIC in %s: %v", name, r)
				result = mcp.NewToolResultError(fmt.Sprintf("Tool %s encountered an internal error", name))
				err = nil // we handled the panic, so no error to return
			}
		}()
		return handler(ctx, req)
	}
}

// withResourceRecovery catches panics in resource handlers
// like having a fire extinguisher in the library - hopefully never needed, but essential
func WithResourceRecovery(uri string, handler server.ResourceHandlerFunc) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) (contents []mcp.ResourceContents, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[RESOURCE] PANIC reading %s: %v", uri, r)
				contents = nil
				err = fmt.Errorf("resource %s encountered an internal error", uri)
			}
		}()
		return handler(ctx, req)
	}
}

// withPromptRecovery catches panics in prompt handlers
// because even conversation templates can sometimes go sideways
func WithPromptRecovery(name string, handler server.PromptHandlerFunc) server.PromptHandlerFunc {
	return func(ctx context.Context, req mcp.GetPromptRequest) (result *mcp.GetPromptResult, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PROMPT] PANIC in %s: %v", name, r)
				result = nil
				err = fmt.Errorf("prompt %s encountered an internal error", name)
			}
		}()
		return handler(ctx, req)
	}
}

// withToolMiddleware combines logging, recovery, and metrics middleware
// this is the full safety package - logging, panic recovery, and performance tracking
func WithToolMiddleware(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	h := WithToolLogging(name, WithToolRecovery(name, handler))
	return WithToolMetrics(name, h)
}

// withResourceMiddleware applies the full middleware stack to resources
// because resources deserve the same level of care as tools
func WithResourceMiddleware(uri string, handler server.ResourceHandlerFunc) server.ResourceHandlerFunc {
	h := WithResourceLogging(uri, WithResourceRecovery(uri, handler))
	return WithResourceMetrics(uri, h)
}

// withPromptMiddleware applies the full middleware stack to prompts
// even conversation templates need proper monitoring
func WithPromptMiddleware(name string, handler server.PromptHandlerFunc) server.PromptHandlerFunc {
	h := WithPromptLogging(name, WithPromptRecovery(name, handler))
	return WithPromptMetrics(name, h)
}
