package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server with all capabilities enabled
	s := server.NewMCPServer(
		"hello-world-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),        // Enable tools
		server.WithResourceCapabilities(true, true), // Enable resources (read, subscribe)
		server.WithPromptCapabilities(true),      // Enable prompts
	)

	// Register a simple tool
	registerGreetTool(s)

	// Register a simple resource
	registerInfoResource(s)

	// Register a simple prompt
	registerGreetingPrompt(s)

	// Start server with stdio transport
	log.Println("Starting Hello World MCP server...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// ============================================================================
// TOOL: Greet someone by name
// ============================================================================

func registerGreetTool(s *server.MCPServer) {
	greetTool := mcp.NewTool(
		"greet",
		mcp.WithDescription("Greet someone by name with a friendly message"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the person to greet"),
		),
		mcp.WithString("language",
			mcp.Description("Language for the greeting (optional)"),
			mcp.Enum("english", "spanish", "french", "german"),
		),
	)

	s.AddTool(greetTool, handleGreet)
}

func handleGreet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameter
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing name: %v", err)), nil
	}

	// Extract optional parameter (with default)
	language := request.GetString("language", "english")

	// Generate greeting based on language
	var greeting string
	switch language {
	case "spanish":
		greeting = fmt.Sprintf("¡Hola, %s! ¿Cómo estás?", name)
	case "french":
		greeting = fmt.Sprintf("Bonjour, %s! Comment allez-vous?", name)
	case "german":
		greeting = fmt.Sprintf("Guten Tag, %s! Wie geht es Ihnen?", name)
	default:
		greeting = fmt.Sprintf("Hello, %s! Nice to meet you!", name)
	}

	// Return result
	return mcp.NewToolResultText(greeting), nil
}

// ============================================================================
// RESOURCE: Server information
// ============================================================================

func registerInfoResource(s *server.MCPServer) {
	infoResource := mcp.NewResource(
		"hello://info",
		"Server Information",
		mcp.WithResourceDescription("Basic information about this MCP server"),
		mcp.WithMIMEType("text/plain"),
	)

	s.AddResource(infoResource, handleInfo)
}

func handleInfo(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Generate server info
	info := fmt.Sprintf(`Hello World MCP Server
======================
Version: 1.0.0
Started: %s
Transport: stdio

This is a simple example MCP server demonstrating:
- Tools: greet people in different languages
- Resources: read server information
- Prompts: generate greeting conversation starters

Try asking me to greet someone!
`, time.Now().Format(time.RFC3339))

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     info,
		},
	}, nil
}

// ============================================================================
// PROMPT: Greeting conversation starter
// ============================================================================

func registerGreetingPrompt(s *server.MCPServer) {
	greetingPrompt := mcp.NewPrompt(
		"greeting_starter",
		mcp.WithPromptDescription("Generate a friendly conversation starter"),
		mcp.WithArgument("topic",
			mcp.ArgumentDescription("Topic for conversation (optional)"),
		),
		mcp.WithArgument("formality",
			mcp.ArgumentDescription("Level of formality"),
		),
	)

	s.AddPrompt(greetingPrompt, handleGreetingPrompt)
}

func handleGreetingPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	// Extract arguments
	args := request.Params.Arguments
	topic := getStringArg(args, "topic", "general chat")
	formality := getStringArg(args, "formality", "casual")

	// Build prompt messages
	systemMessage := "You are a friendly assistant helping to start conversations."

	var userMessage string
	if formality == "formal" {
		userMessage = fmt.Sprintf("Generate a formal, professional conversation starter about: %s", topic)
	} else {
		userMessage = fmt.Sprintf("Generate a casual, friendly conversation starter about: %s", topic)
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Greeting prompt for topic: %s", topic),
		Messages: []mcp.PromptMessage{
			{
				Role: "system",
				Content: mcp.NewTextContent(systemMessage),
			},
			{
				Role: "user",
				Content: mcp.NewTextContent(userMessage),
			},
		},
	}, nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// getStringArg safely extracts a string argument with a default value
func getStringArg(args map[string]string, key, defaultValue string) string {
	if val, ok := args[key]; ok {
		return val
	}
	return defaultValue
}
