package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolMiddleware wraps tool handlers with error handling and logging
func WithToolLogging(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()
		log.Printf("[TOOL] Starting: %s with args: %v", name, req.Params.Arguments)

		result, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			log.Printf("[TOOL] Error in %s (took %v): %v", name, duration, err)
			// Return a safe error message to the client
			return mcp.NewToolResultError(fmt.Sprintf("Tool %s failed: %v", name, err)), nil
		}

		log.Printf("[TOOL] Completed: %s (took %v)", name, duration)
		return result, nil
	}
}

// ResourceMiddleware wraps resource handlers with error handling and logging
func WithResourceLogging(uri string, handler server.ResourceHandlerFunc) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		start := time.Now()
		log.Printf("[RESOURCE] Reading: %s", uri)

		contents, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			log.Printf("[RESOURCE] Error reading %s (took %v): %v", uri, duration, err)
			return nil, fmt.Errorf("failed to read resource %s: %w", uri, err)
		}

		log.Printf("[RESOURCE] Successfully read: %s (took %v, %d items)", uri, duration, len(contents))
		return contents, nil
	}
}

// PromptMiddleware wraps prompt handlers with error handling and logging
func WithPromptLogging(name string, handler server.PromptHandlerFunc) server.PromptHandlerFunc {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		start := time.Now()
		log.Printf("[PROMPT] Getting: %s with args: %v", name, req.Params.Arguments)

		result, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			log.Printf("[PROMPT] Error in %s (took %v): %v", name, duration, err)
			return nil, fmt.Errorf("failed to get prompt %s: %w", name, err)
		}

		log.Printf("[PROMPT] Completed: %s (took %v)", name, duration)
		return result, nil
	}
}

// Recovery middleware to catch panics
func WithToolRecovery(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[TOOL] PANIC in %s: %v", name, r)
				result = mcp.NewToolResultError(fmt.Sprintf("Tool %s encountered an internal error", name))
				err = nil
			}
		}()
		return handler(ctx, req)
	}
}

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

// Combine logging, recovery, and metrics middleware
func WithToolMiddleware(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	h := WithToolLogging(name, WithToolRecovery(name, handler))
	return WithToolMetrics(name, h)
}

func WithResourceMiddleware(uri string, handler server.ResourceHandlerFunc) server.ResourceHandlerFunc {
	h := WithResourceLogging(uri, WithResourceRecovery(uri, handler))
	return WithResourceMetrics(uri, h)
}

func WithPromptMiddleware(name string, handler server.PromptHandlerFunc) server.PromptHandlerFunc {
	h := WithPromptLogging(name, WithPromptRecovery(name, handler))
	return WithPromptMetrics(name, h)
}
