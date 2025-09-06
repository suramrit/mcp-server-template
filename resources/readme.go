package resources

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// readmeResource gives AI access to our README file
// think of this as our helpful librarian that fetches books on demand
type ReadmeResource struct {
	// no state needed - we're just a simple file reader
}

// newReadmeResource creates our file-reading resource
// every project needs documentation, and every AI needs access to it!
func NewReadmeResource() *ReadmeResource {
	return &ReadmeResource{}
}

// getResource defines what this resource looks like to MCP clients
// this is like putting a label on a library book
func (r *ReadmeResource) GetResource() mcp.Resource {
	return mcp.NewResource(
		"file://README.md",           // unique URI - like a library catalog number
		"Local README",               // human-readable description
		mcp.WithMIMEType("text/plain")) // tell clients what kind of content this is
}

// getHandler returns the function that actually reads the file
// this is our librarian in action - fetching the requested book
func (r *ReadmeResource) GetHandler() server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// try to read the file from our resources directory
		// note: this path is relative to where the server runs
		b, err := os.ReadFile("resources/static/README.md")
		if err != nil {
			// if file doesn't exist or can't be read, let the caller know
			return nil, err
		}
		
		// package up the file contents in the format MCP expects
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",    // internal URI for this content
				MIMEType: "text/markdown",    // what kind of content this actually is
				Text:     string(b),          // the actual file contents
			},
		}, nil
	}
}
