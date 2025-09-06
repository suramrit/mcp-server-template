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

// main is where the magic begins!
// this is the entry point for our MCP server that will give AI superpowers
func main() {
	// create a log file because stdout/stderr get hijacked by MCP protocol
	// think of this as our server's diary - it'll tell us everything that happens!
	logFile, err := os.OpenFile("mcp-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err) // if we can't log, we can't debug - bail out!
	}
	defer logFile.Close() // always clean up after yourself, mom would be proud

	// redirect all our chatty log messages to the file
	log.SetOutput(logFile)
	log.Println("=== MCP Server Starting ===") // and... action!

	// build our MCP server - this is the foundation everything sits on
	srv := server.NewMCPServer(
		"hello-mcp",                      // server name - keep it friendly!
		"0.1.0",                         // version - we're just getting started
		server.WithToolCapabilities(false), // we'll handle tool capabilities ourselves
		server.WithLogging(),               // enable the SDK's internal logging too
	)

	// time to set up the three-ring circus of MCP capabilities!
	log.Println("Registering tools, resources, and prompts...")
	
	// tools: let AI DO things (like our friendly echo)
	registerTools(srv, tools.NewEchoTool())
	
	// resources: give AI access to data (like our README file)
	registerResources(srv, resources.NewReadmeResource())
	
	// prompts: provide AI with conversation templates
	registerPrompts(srv, prompts.NewGreetingPrompt())

	log.Println("Starting stdio server...")
	// launch! This blocks forever, listening for JSON-RPC over stdin/stdout
	// the MCP client (like Claude Desktop) will spawn us and talk to us here
	if err := server.ServeStdio(srv); err != nil {
		log.Fatal(err) // if the server dies, we die with it
	}
}

// registerTools wraps our tools with middleware and registers them
// think of this as putting on safety gear before using power tools!
func registerTools(srv *server.MCPServer, tools ...tools.Tool) {
	for _, tool := range tools {
		toolDef := tool.GetTool()   // get the tool's definition (name, params, etc.)
		handler := tool.GetHandler() // get the actual function that does the work

		// wrap with middleware: logging, metrics, panic recovery
		// it's like having a safety net for our trapeze artists!
		wrappedHandler := middleware.WithToolMiddleware(toolDef.Name, handler)

		// register with the server - now AI can call this tool!
		srv.AddTool(toolDef, wrappedHandler)
	}
}

// registerResources gives AI access to data sources
// like giving AI a library card!
func registerResources(srv *server.MCPServer, resources ...resources.Resource) {
	for _, resource := range resources {
		resourceDef := resource.GetResource() // what is this resource?
		handler := resource.GetHandler()      // how do we read it?

		// same middleware magic - safety first!
		wrappedHandler := middleware.WithResourceMiddleware(resourceDef.URI, handler)

		// now AI can ask for this data whenever it needs it
		srv.AddResource(resourceDef, wrappedHandler)
	}
}

// registerPrompts sets up conversation templates for AI to use
// think of these as conversation starters or script templates
func registerPrompts(srv *server.MCPServer, prompts ...prompts.Prompt) {
	for _, prompt := range prompts {
		promptDef := prompt.GetPrompt() // what kind of prompt is this?
		handler := prompt.GetHandler()  // how do we generate the template?

		// middleware wrapping - because we love consistency!
		wrappedHandler := middleware.WithPromptMiddleware(promptDef.Name, handler)

		// register the prompt so AI can use our templates
		srv.AddPrompt(promptDef, wrappedHandler)
	}
}
