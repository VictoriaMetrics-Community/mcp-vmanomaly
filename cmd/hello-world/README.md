# Hello World MCP Server

A simple example MCP (Model Context Protocol) server demonstrating the three main MCP features: **Tools**, **Resources**, and **Prompts**.

## What This Demonstrates

This server shows you how to:
- ✅ Create an MCP server with stdio transport
- ✅ Register a **Tool** with parameters (required and optional)
- ✅ Register a **Resource** that provides readable content
- ✅ Register a **Prompt** template with arguments

## Building

```bash
# From project root
go build -o bin/hello-world ./cmd/hello-world

# Or use make (if you add it to Makefile)
make build
```

## Running

The server uses **stdio transport**, which means it communicates via standard input/output:

```bash
./bin/hello-world
```

The server will start and wait for MCP protocol messages on stdin.

## Testing with Claude Desktop

Add this to your Claude Desktop configuration:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "hello-world": {
      "command": "/absolute/path/to/bin/hello-world"
    }
  }
}
```

After adding, restart Claude Desktop.

## Available Features

### Tool: `greet`

Greet someone in different languages.

**Parameters**:
- `name` (required): Name of the person to greet
- `language` (optional): Language - `english`, `spanish`, `french`, `german` (default: `english`)

**Example**:
```
You: Greet Alice in Spanish
Claude: Uses the greet tool with name="Alice", language="spanish"
Result: ¡Hola, Alice! ¿Cómo estás?
```

### Resource: `hello://info`

Read server information and capabilities.

**Example**:
```
You: Read the hello://info resource
Claude: Fetches the resource
Result: Shows server version, features, and description
```

### Prompt: `greeting_starter`

Generate conversation starters.

**Arguments**:
- `topic` (optional): Topic for conversation (default: `general chat`)
- `formality` (optional): `casual` or `formal` (default: `casual`)

**Example**:
```
You: Use the greeting_starter prompt with topic="technology"
Claude: Uses the prompt template
Result: Generates a friendly conversation starter about technology
```

## Code Structure

```go
// main.go structure:
main()                      // Server initialization
├── registerGreetTool()     // Tool: greet with language options
│   └── handleGreet()       // Tool handler function
├── registerInfoResource()  // Resource: server info
│   └── handleInfo()        // Resource handler function
└── registerGreetingPrompt() // Prompt: conversation starter
    └── handleGreetingPrompt() // Prompt handler function
```

## Key Learning Points

### 1. Server Creation
```go
s := server.NewMCPServer(
    "hello-world-mcp",
    "1.0.0",
    server.WithToolCapabilities(true),        // Enable tools
    server.WithResourceCapabilities(true, true), // Enable resources
    server.WithPromptCapabilities(true),      // Enable prompts
)
```

### 2. Tool Registration
```go
tool := mcp.NewTool("greet",
    mcp.WithDescription("Greet someone by name"),
    mcp.WithString("name", mcp.Required()),
    mcp.WithString("language", mcp.Enum("english", "spanish")),
)
s.AddTool(tool, handleGreet)
```

### 3. Resource Registration
```go
resource := mcp.NewResource("hello://info", "Server Information",
    mcp.WithMIMEType("text/plain"),
)
s.AddResource(resource, handleInfo)
```

### 4. Prompt Registration
```go
prompt := mcp.NewPrompt("greeting_starter",
    mcp.WithPromptDescription("Generate a conversation starter"),
    mcp.WithArgument("topic", mcp.ArgumentDescription("Topic")),
)
s.AddPrompt(prompt, handleGreetingPrompt)
```

### 5. Starting Server
```go
// stdio transport (default for MCP)
server.ServeStdio(s)
```

## Handler Signatures

**Tool Handler**:
```go
func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
```

**Resource Handler**:
```go
func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
```

**Prompt Handler**:
```go
func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error)
```

## Next Steps

After understanding this example:
1. Modify the tool to add more languages
2. Add another resource (e.g., `hello://stats`)
3. Create a new tool that does something useful
4. Try different parameter types (numbers, booleans, arrays)
5. Look at the main vmanomaly server implementation

## Troubleshooting

**Server won't start in Claude Desktop**:
- Check the path in config is absolute
- Verify binary has execute permissions: `chmod +x bin/hello-world`
- Check Claude Desktop logs for errors

**Changes not showing**:
- Restart Claude Desktop after config changes
- Rebuild the binary after code changes

**Testing without Claude Desktop**:
```bash
# Send MCP initialize message (advanced)
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"0.1.0","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./bin/hello-world
```

## References

- [MCP Documentation](https://modelcontextprotocol.io/)
- [mcp-go Library](https://github.com/mark3labs/mcp-go)
- [mcp-go Docs](https://mcp-go.dev/)
