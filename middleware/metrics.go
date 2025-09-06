package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Metrics struct {
	mu sync.RWMutex

	ToolCalls     map[string]int64
	ToolErrors    map[string]int64
	ToolDurations map[string]time.Duration

	ResourceReads     map[string]int64
	ResourceErrors    map[string]int64
	ResourceDurations map[string]time.Duration

	PromptGets      map[string]int64
	PromptErrors    map[string]int64
	PromptDurations map[string]time.Duration
}

var GlobalMetrics = &Metrics{
	ToolCalls:         make(map[string]int64),
	ToolErrors:        make(map[string]int64),
	ToolDurations:     make(map[string]time.Duration),
	ResourceReads:     make(map[string]int64),
	ResourceErrors:    make(map[string]int64),
	ResourceDurations: make(map[string]time.Duration),
	PromptGets:        make(map[string]int64),
	PromptErrors:      make(map[string]int64),
	PromptDurations:   make(map[string]time.Duration),
}

func WithToolMetrics(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()

		GlobalMetrics.mu.Lock()
		GlobalMetrics.ToolCalls[name]++
		GlobalMetrics.mu.Unlock()

		result, err := handler(ctx, req)

		duration := time.Since(start)
		GlobalMetrics.mu.Lock()
		GlobalMetrics.ToolDurations[name] += duration
		if err != nil {
			GlobalMetrics.ToolErrors[name]++
		}
		GlobalMetrics.mu.Unlock()

		return result, err
	}
}

func WithResourceMetrics(uri string, handler server.ResourceHandlerFunc) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		start := time.Now()

		GlobalMetrics.mu.Lock()
		GlobalMetrics.ResourceReads[uri]++
		GlobalMetrics.mu.Unlock()

		contents, err := handler(ctx, req)

		duration := time.Since(start)
		GlobalMetrics.mu.Lock()
		GlobalMetrics.ResourceDurations[uri] += duration
		if err != nil {
			GlobalMetrics.ResourceErrors[uri]++
		}
		GlobalMetrics.mu.Unlock()

		return contents, err
	}
}

func WithPromptMetrics(name string, handler server.PromptHandlerFunc) server.PromptHandlerFunc {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		start := time.Now()

		GlobalMetrics.mu.Lock()
		GlobalMetrics.PromptGets[name]++
		GlobalMetrics.mu.Unlock()

		result, err := handler(ctx, req)

		duration := time.Since(start)
		GlobalMetrics.mu.Lock()
		GlobalMetrics.PromptDurations[name] += duration
		if err != nil {
			GlobalMetrics.PromptErrors[name]++
		}
		GlobalMetrics.mu.Unlock()

		return result, err
	}
}

func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"tool_calls":         m.ToolCalls,
		"tool_errors":        m.ToolErrors,
		"tool_durations":     m.ToolDurations,
		"resource_reads":     m.ResourceReads,
		"resource_errors":    m.ResourceErrors,
		"resource_durations": m.ResourceDurations,
		"prompt_gets":        m.PromptGets,
		"prompt_errors":      m.PromptErrors,
		"prompt_durations":   m.PromptDurations,
	}
}
