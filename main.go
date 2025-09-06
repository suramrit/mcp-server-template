package main

import (
	"log"
	"os"

	server "github.com/mark3labs/mcp-go/server"
	"github.com/suramrit/hello-mcp/middleware"
	"github.com/suramrit/hello-mcp/prompts"
	"github.com/suramrit/hello-mcp/resources"
	"github.com/suramrit/hello-mcp/tools"
)

func main() {
	// Create a log file to see what's happening
	logFile, err := os.OpenFile("mcp-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Set log output to file
	log.SetOutput(logFile)
	log.Println("=== MCP Server Starting ===")

	// stdio transport so desktop MCP clients can spawn us as a child process
	srv := server.NewMCPServer(
		"hello-mcp",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithLogging(),
	)

	log.Println("Registering tools, resources, and prompts...")
	registerTools(srv, tools.NewEchoTool())
	registerResources(srv, resources.NewReadmeResource())
	registerPrompts(srv, prompts.NewGreetingPrompt())

	log.Println("Starting stdio server...")
	// Serve over stdio; blocks
	if err := server.ServeStdio(srv); err != nil {
		log.Fatal(err)
	}
}

func registerTools(srv *server.MCPServer, tools ...tools.Tool) {
	for _, tool := range tools {
		toolDef := tool.GetTool()
		handler := tool.GetHandler()

		// Wrap handler with middleware
		wrappedHandler := middleware.WithToolMiddleware(toolDef.Name, handler)

		srv.AddTool(toolDef, wrappedHandler)
	}
}

func registerResources(srv *server.MCPServer, resources ...resources.Resource) {
	for _, resource := range resources {
		resourceDef := resource.GetResource()
		handler := resource.GetHandler()

		// Wrap handler with middleware
		wrappedHandler := middleware.WithResourceMiddleware(resourceDef.URI, handler)

		srv.AddResource(resourceDef, wrappedHandler)
	}
}

func registerPrompts(srv *server.MCPServer, prompts ...prompts.Prompt) {
	for _, prompt := range prompts {
		promptDef := prompt.GetPrompt()
		handler := prompt.GetHandler()

		// Wrap handler with middleware
		wrappedHandler := middleware.WithPromptMiddleware(promptDef.Name, handler)

		srv.AddPrompt(promptDef, wrappedHandler)
	}
}
