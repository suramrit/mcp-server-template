package resources

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ReadmeResource struct {
}

func NewReadmeResource() *ReadmeResource {
	return &ReadmeResource{}
}

func (r *ReadmeResource) GetResource() mcp.Resource {
	return mcp.NewResource("file://README.md", "Local README", mcp.WithMIMEType("text/plain"))
}

func (r *ReadmeResource) GetHandler() server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		b, err := os.ReadFile("resources/static/README.md")
		if err != nil {
			return nil, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(b),
			},
		}, nil
	}
}
