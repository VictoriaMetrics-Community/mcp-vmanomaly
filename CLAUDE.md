# MCP vmanomaly - Claude Memory File

This file contains important context and memory for Claude when working on the mcp-vmanomaly project.

## Project Overview

**Project**: MCP Server for vmanomaly
**Language**: Go 1.25
**Purpose**: Model Context Protocol server to enable AI assistants to interact with vmanomaly REST API

vmanomaly is VictoriaMetrics' anomaly detection application built on top of VictoriaMetrics and VictoriaLogs. This MCP server acts as a bridge between AI assistants (like Claude) and the vmanomaly REST API.

## Developer Context

**Developer Background**:
- New to Go (primary experience: 10+ years Python, 1+ years TypeScript)
- Familiar with REST APIs, HTTP clients, and backend development
- Working with VictoriaMetrics ecosystem

**Learning Resources for Go**:
- Go is statically typed with explicit error handling (no exceptions)
- Package management via go.mod (similar to package.json or requirements.txt)
- Interfaces are implicit (no explicit "implements" keyword)
- Error handling pattern: `if err != nil { return err }`
- Contexts are passed as first parameter for cancellation/timeouts

## Architecture

### Directory Structure

```
mcp-vmanomaly/
├── cmd/mcp-vmanomaly/          # Main application entry point
│   └── main.go                 # Server initialization and startup
├── internal/                   # Private application code
│   ├── vmanomaly/             # vmanomaly API client package
│   │   └── client.go          # HTTP client for vmanomaly REST API
│   └── tools/                 # MCP tool definitions
│       └── tools.go           # Tool registration and handlers
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
├── Makefile                   # Build and development tasks
├── README.md                  # User-facing documentation
└── CLAUDE.md                  # This file - Claude's memory
```

### Key Dependencies

1. **github.com/mark3labs/mcp-go** - MCP protocol implementation for Go
   - Provides server framework, tool definitions, and protocol handling
   - Main types: `mcp.Tool`, `server.MCPServer`, `mcp.CallToolRequest`

2. **github.com/tmc/langchaingo** - LangChain for Go
   - Provides AI/LLM integrations and utilities
   - May be used for advanced AI features in the future

3. **net/http** (stdlib) - HTTP client for API calls
   - Used in `internal/vmanomaly/client.go` to call vmanomaly REST API

### Design Patterns

**Client Pattern** (`internal/vmanomaly/client.go`):
- Encapsulates all HTTP communication with vmanomaly
- Uses Go's `context.Context` for request cancellation and timeouts
- Centralizes authentication (Bearer token in Authorization header)
- Generic `doRequest()` method for all API calls

**Tool Registration Pattern** (`internal/tools/tools.go`):
- Each MCP tool is defined with `mcp.NewTool()`
- Handler functions follow signature: `func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)`
- Tools are registered in `RegisterTools()` function
- Handlers call client methods and format responses

**Error Handling**:
- Always check errors: `if err != nil { return err }`
- Use `fmt.Errorf()` with `%w` for error wrapping
- Return `mcp.NewToolResultError()` for tool errors
- Log fatal errors in main.go with `log.Fatal()`

## Configuration

**Environment Variables**:
- `VMANOMALY_ENDPOINT` (required): vmanomaly server URL (e.g., `http://localhost:8490`)
- `VMANOMALY_API_TOKEN` (optional): Authentication token for vmanomaly API

**Transport**: MCP uses stdio transport by default (communication via stdin/stdout)

## vmanomaly REST API

**Important**: The actual vmanomaly REST API endpoints need to be documented/discovered.

**Known/Assumed Endpoints**:
- `/health` - Health check (implemented)
- TBD: Model management endpoints
- TBD: Anomaly detection endpoints
- TBD: Metrics endpoints

**TODO**: Reference vmanomaly API documentation to implement additional endpoints.

## Development Workflow

### Common Tasks

```bash
# Install/update dependencies
make install

# Build the binary
make build

# Run the server (needs VMANOMALY_ENDPOINT)
VMANOMALY_ENDPOINT=http://localhost:8490 make run

# Run tests
make test

# Format code (Go convention: use gofmt)
make fmt

# Lint/vet code
make lint

# Clean build artifacts
make clean
```

### Adding a New Tool

**Steps**:
1. **Add API method** to `internal/vmanomaly/client.go`
   ```go
   func (c *Client) MethodName(ctx context.Context, params ...) (Result, error) {
       return c.doRequest(ctx, http.MethodGET, "/endpoint", nil)
   }
   ```

2. **Create tool definition** in `internal/tools/tools.go`
   ```go
   tool := mcp.NewTool("tool_name",
       mcp.WithDescription("Tool description"),
       mcp.WithString("param_name", mcp.Required(true)),
   )
   ```

3. **Register handler** in `RegisterTools()`
   ```go
   s.AddTool(tool, handlerFunc(client))
   ```

4. **Implement handler**
   ```go
   func handlerFunc(client *vmanomaly.Client) server.ToolHandler {
       return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
           // Call client method
           // Format and return result
       }
   }
   ```

5. **Update documentation** in README.md

### Testing

- Use `go test` for unit tests
- Test files should be named `*_test.go`
- Table-driven tests are idiomatic in Go
- Use `context.Background()` for test contexts

### Go-Specific Gotchas for Python/TypeScript Developers

1. **No implicit type conversion** - must explicitly convert types
2. **Multiple return values** - common pattern: `(result, error)`
3. **Exported vs unexported** - uppercase = public, lowercase = private
4. **Nil vs undefined** - zero value for pointers/interfaces is `nil`
5. **Defer** - cleanup code (like Python's context managers): `defer resp.Body.Close()`
6. **Channels and goroutines** - Go's concurrency primitives (not needed yet)

## Reference Implementation

**Inspiration**: [mcp-victoriametrics](https://github.com/VictoriaMetrics-Community/mcp-victoriametrics)
- Similar MCP server for VictoriaMetrics
- Good reference for tool patterns and structure
- Uses similar architecture and dependencies

## Future Enhancements

**Planned Tools** (from README):
- `list_models` - List available anomaly detection models
- `run_detection` - Trigger anomaly detection
- `get_anomalies` - Retrieve detected anomalies
- `get_model_metrics` - Get model performance metrics
- `query_logs` - Query VictoriaLogs for anomaly data

**Potential Features**:
- Streaming support for long-running detection jobs
- Caching for frequently accessed data
- Resource providers (MCP resources for model configs)
- Prompts for common anomaly detection workflows

## Important Notes

- **Stateless design**: MCP server should be stateless, all state in vmanomaly
- **Error messages**: Should be clear and actionable for AI assistant users
- **JSON formatting**: Use `json.MarshalIndent()` for readable tool responses
- **Context propagation**: Always pass `context.Context` through call chain
- **Documentation**: Keep README.md in sync with new tools/features

## Quick Reference

**Build the project**: `make build`
**Run tests**: `make test`
**View all make targets**: `make help`
**Format code**: `make fmt`

**MCP protocol docs**: https://modelcontextprotocol.io/
**vmanomaly docs**: https://docs.victoriametrics.com/anomaly-detection/
**mcp-go examples**: https://github.com/mark3labs/mcp-go/tree/main/examples

---

**Last Updated**: 2025-11-09
**Status**: Initial project setup complete, basic health check tool implemented
