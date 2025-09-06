package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// metrics holds all our performance data
// think of this as our server's fitness tracker - it knows everything about performance
type Metrics struct {
	mu sync.RWMutex // protects all the maps below from concurrent access chaos

	// tool metrics - how are our tools performing?
	ToolCalls     map[string]int64        // how many times each tool was called
	ToolErrors    map[string]int64        // how many times each tool failed
	ToolDurations map[string]time.Duration // total time spent in each tool

	// resource metrics - how's our data access doing?
	ResourceReads     map[string]int64        // how many times each resource was read
	ResourceErrors    map[string]int64        // how many read failures we've had
	ResourceDurations map[string]time.Duration // time spent reading resources

	// prompt metrics - are our conversation templates popular?
	PromptGets      map[string]int64        // how many times each prompt was requested
	PromptErrors    map[string]int64        // prompt generation failures (shouldn't be many!)
	PromptDurations map[string]time.Duration // time spent generating prompts
}

// globalMetrics is our singleton metrics collector
// there's only one performance tracker per server - we're not running a democracy here
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

// withToolMetrics wraps tool handlers to collect performance data
// this is like having a stopwatch and counter for every tool
func WithToolMetrics(name string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()

		// increment the call counter - thread-safely!
		GlobalMetrics.mu.Lock()
		GlobalMetrics.ToolCalls[name]++
		GlobalMetrics.mu.Unlock()

		// do the actual work
		result, err := handler(ctx, req)

		// record how long it took and whether it succeeded
		duration := time.Since(start)
		GlobalMetrics.mu.Lock()
		GlobalMetrics.ToolDurations[name] += duration
		if err != nil {
			GlobalMetrics.ToolErrors[name]++ // another one bites the dust
		}
		GlobalMetrics.mu.Unlock()

		return result, err
	}
}

// withResourceMetrics tracks resource reading performance
// because knowing how fast we can read data is important for scaling
func WithResourceMetrics(uri string, handler server.ResourceHandlerFunc) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		start := time.Now()

		// count this read attempt
		GlobalMetrics.mu.Lock()
		GlobalMetrics.ResourceReads[uri]++
		GlobalMetrics.mu.Unlock()

		// try to read the resource
		contents, err := handler(ctx, req)

		// record the results
		duration := time.Since(start)
		GlobalMetrics.mu.Lock()
		GlobalMetrics.ResourceDurations[uri] += duration
		if err != nil {
			GlobalMetrics.ResourceErrors[uri]++ // file not found? permission denied? we'll know!
		}
		GlobalMetrics.mu.Unlock()

		return contents, err
	}
}

// withPromptMetrics tracks prompt generation performance
// even conversation templates need performance monitoring
func WithPromptMetrics(name string, handler server.PromptHandlerFunc) server.PromptHandlerFunc {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		start := time.Now()

		// count this prompt request
		GlobalMetrics.mu.Lock()
		GlobalMetrics.PromptGets[name]++
		GlobalMetrics.mu.Unlock()

		// generate the prompt
		result, err := handler(ctx, req)

		// record the performance data
		duration := time.Since(start)
		GlobalMetrics.mu.Lock()
		GlobalMetrics.PromptDurations[name] += duration
		if err != nil {
			GlobalMetrics.PromptErrors[name]++ // this really shouldn't happen often
		}
		GlobalMetrics.mu.Unlock()

		return result, err
	}
}

// getStats returns a snapshot of all our performance metrics
// this is like getting a report card for your server's performance
func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// return everything we've collected - tools, resources, and prompts
	return map[string]interface{}{
		"tool_calls":         m.ToolCalls,      // how busy are our tools?
		"tool_errors":        m.ToolErrors,     // how reliable are they?
		"tool_durations":     m.ToolDurations,  // how fast are they?
		"resource_reads":     m.ResourceReads,  // how much data are we serving?
		"resource_errors":    m.ResourceErrors, // any file access problems?
		"resource_durations": m.ResourceDurations, // how fast is our I/O?
		"prompt_gets":        m.PromptGets,     // are our templates popular?
		"prompt_errors":      m.PromptErrors,   // any template generation issues?
		"prompt_durations":   m.PromptDurations, // how fast can we generate templates?
	}
}
