package resources

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// resource interface defines what every resource must provide
// like a library card catalog - every resource needs these basic details
type Resource interface {
	GetResource() mcp.Resource                // what resource are you? (URI, description, type)
	GetHandler() server.ResourceHandlerFunc   // how do we read you? (the actual reading logic)
}
