package resources

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Resource interface {
	GetResource() mcp.Resource
	GetHandler() server.ResourceHandlerFunc
}
