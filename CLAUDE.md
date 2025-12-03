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
│   │   ├── client.go          # HTTP client for vmanomaly REST API
│   │   └── types.go           # Request/response type definitions
│   ├── resources/             # MCP resources (embedded documentation)
│   │   ├── docs.go            # Resource registration and search
│   │   ├── utils.go           # Markdown processing utilities
│   │   ├── docs_test.go       # Resource tests
│   │   └── docs/              # Embedded documentation files (~5.4 MB)
│   └── tools/                 # MCP tool definitions
│       ├── tools.go           # Tool registration (entry point)
│       ├── models.go          # Model configuration tools (3)
│       ├── tasks.go           # Anomaly detection task tools (5) [WIP]
│       ├── query.go           # Query tool (1)
│       ├── info.go            # Info/utility tools (2)
│       ├── config.go          # Configuration tools (1)
│       ├── compatibility.go   # Compatibility check tool (1)
│       ├── alerts.go          # Alert rule generation tools (1)
│       └── docs.go            # Documentation search tool (1)
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

2. **github.com/blevesearch/bleve/v2** - Full-text search engine for Go
   - Used for documentation search indexing and querying
   - Provides fuzzy matching and relevance scoring
   - In-memory index for fast documentation lookups

3. **github.com/tmc/langchaingo** - LangChain for Go
   - Provides text splitters for markdown chunking
   - Used to split large docs into 65KB chunks with 8KB overlap
   - Preserves heading hierarchy and code blocks

4. **net/http** (stdlib) - HTTP client for API calls
   - Used in `internal/vmanomaly/client.go` to call vmanomaly REST API

### Design Patterns

**Client Pattern** (`internal/vmanomaly/client.go`):
- Encapsulates all HTTP communication with vmanomaly
- Uses Go's `context.Context` for request cancellation and timeouts
- Generic `doRequest()` method for all API calls

**Tool Registration Pattern** (`internal/tools/*.go`):
- Tools use struct-based schema definitions with `mcp.WithInputSchema[T]()`
- Handler functions use `mcp.NewTypedToolHandler` pattern
- Handler signature: `func(ctx context.Context, req mcp.CallToolRequest, args T) (*mcp.CallToolResult, error)`
- Args structs use `jsonschema_description` tags for parameter documentation
- Tools are organized in separate files: models.go, tasks.go, query.go, info.go, config.go, compatibility.go, alerts.go, docs.go
- Main registration in `RegisterTools()` calls sub-registration functions

**Error Handling**:
- Always check errors: `if err != nil { return err }`
- Use `fmt.Errorf()` with `%w` for error wrapping
- Return `mcp.NewToolResultError()` for tool errors
- Log fatal errors in main.go with `log.Fatal()`

**Resources Pattern** (`internal/resources/*.go`):
- Embedded documentation using `//go:embed docs` directive
- Three-tier caching strategy:
  1. Bleve search index (full-text search with fuzzy matching)
  2. Resources map (resource metadata by URI)
  3. Contents map (actual content by URI)
- Markdown chunking with langchaingo's MarkdownTextSplitter:
  - 65KB chunk size, 8KB overlap
  - Preserves heading hierarchy, code blocks, table rows
  - Front matter extraction (YAML title → H1 heading)
- URI format: `docs://anomaly-detection/path/file.md#chunk_num`
- Single handler pattern: one `docResourcesHandler()` serves all resources
- Search returns embedded resources via `mcp.EmbeddedResource`

## Configuration

**Environment Variables**:
- `VMANOMALY_ENDPOINT` (required): vmanomaly server URL (e.g., `http://localhost:8490`)

**Transport**: MCP uses stdio transport by default (communication via stdin/stdout)

## vmanomaly REST API

**Fully Implemented** - All major endpoints from OpenAPI spec (http://localhost:8490/docs):

### Model Configuration Endpoints
- `GET /api/v1/models` - List available model types
- `GET /api/v1/model/schema?model_class={class}` - Get JSON schema for model
- `POST /api/v1/model/validate` - Validate model configuration
- `GET /api/vmanomaly/config.yaml` - Generate complete YAML config

### Anomaly Detection Task Endpoints
- `POST /api/v1/anomaly_detection/tasks` - Create detection task
- `GET /api/v1/anomaly_detection/tasks?limit={n}&status={status}` - List tasks
- `GET /api/v1/anomaly_detection/tasks/{task_id}` - Get task status
- `DELETE /api/v1/anomaly_detection/tasks/{task_id}` - Cancel task
- `GET /api/v1/anomaly_detection/limits` - Get system capacity/limits

### Query & Utility Endpoints
- `POST /api/v1/query` - Execute PromQL/LogsQL query
- `GET /api/v1/status/buildinfo` - Get build information
- `GET /health` - Health check

### Alert Rule Endpoints
- `GET /api/vmanomaly/example-alert-rule.yaml` - Generate VMAlert rule configuration

**Available Model Types**:
zscore, prophet, mad, holtwinters, std, rolling_quantile, isolation_forest_univariate, mad_online, zscore_online, quantile_online, auto

**Response Formats**: All endpoints return JSON (except /api/vmanomaly/config.yaml which returns YAML)

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

## MCP Struct-Based Schema Pattern (IMPORTANT!)

**We use the modern struct-based schema approach**, not the builder pattern:

```go
// 1. Define args struct with json and jsonschema_description tags
type ToolNameArgs struct {
    RequiredParam string  `json:"required_param" jsonschema_description:"Parameter description"`
    OptionalParam string  `json:"optional_param,omitempty" jsonschema_description:"Optional parameter"`
}

// 2. Register tool with WithInputSchema[T]()
tool := mcp.NewTool("tool_name",
    mcp.WithDescription("Tool description"),
    mcp.WithInputSchema[ToolNameArgs](),  // NOT WithSchema(ToolNameArgs{})
)

// 3. Use NewTypedToolHandler wrapper
s.AddTool(tool, mcp.NewTypedToolHandler(handleToolName(client)))

// 4. Handler signature includes args as parameter
func handleToolName(client *Client) func(ctx context.Context, req mcp.CallToolRequest, args ToolNameArgs) (*mcp.CallToolResult, error) {
    return func(ctx context.Context, req mcp.CallToolRequest, args ToolNameArgs) (*mcp.CallToolResult, error) {
        // args.RequiredParam is already parsed and typed!
        // No need for req.UnmarshalArguments()
    }
}
```

**Key Differences from Builder Pattern**:
- ✅ Use `mcp.WithInputSchema[T]()` (generic)
- ✅ Use `mcp.NewTypedToolHandler(handler)` wrapper
- ✅ Handler gets `args T` as third parameter
- ✅ Use `jsonschema_description` tag (with underscore!)
- ❌ NOT `mcp.WithString()`, `mcp.WithNumber()`, etc.
- ❌ NOT `mcp.WithSchema(T{})` (non-generic)
- ❌ NOT `req.UnmarshalArguments(&args)`

## Reference Implementation

**Inspiration**: [mcp-victoriametrics](https://github.com/VictoriaMetrics-Community/mcp-victoriametrics)
- Similar MCP server for VictoriaMetrics
- Good reference for tool patterns and structure
- Uses similar architecture and dependencies

## Implementation Status

**✅ COMPLETED** - All major features implemented (13 tools + resources):

### Model Configuration (4 tools)
- ✅ `list_models` - List available model types
- ✅ `get_model_schema` - Get JSON schema for model type
- ✅ `validate_model_config` - Validate model configuration
- ✅ `generate_config` - Generate complete YAML config

### Anomaly Detection Tasks (5 tools)
- ✅ `create_detection_task` - Create and start detection task
- ✅ `get_task_status` - Get task status with progress
- ✅ `list_tasks` - List tasks with filtering
- ✅ `cancel_task` - Cancel running task
- ✅ `get_detection_limits` - Get system capacity

### Query & Utility (3 tools)
- ✅ `query_metrics` - Execute PromQL/LogsQL queries
- ✅ `get_buildinfo` - Get build information
- ✅ `health_check` - Health check endpoint

### Documentation Search (1 tool)
- ✅ `search_docs` - Full-text search across vmanomaly documentation

### MCP Resources
- ✅ Embedded documentation (21 MD files, 30+ images, ~5.4 MB)
- ✅ Bleve full-text search index with fuzzy matching
- ✅ Markdown chunking (65KB chunks, 8KB overlap)
- ✅ ~50+ documentation resources (chunked)
- ✅ URI format: `docs://anomaly-detection/path/file.md#chunk_num`

## Future Enhancements (Optional)

**Potential Features**:
- Additional MCP Resources:
  - `vmanomaly://models` - Cached list of model types
  - `vmanomaly://tasks` - Current running tasks overview
  - `vmanomaly://limits` - System capacity as resource
- MCP Prompts for common workflows:
  - "Create anomaly detection setup" - Guided workflow
  - "Debug failed task" - Troubleshooting assistant
  - "Explain anomaly score" - Interactive documentation
- Streaming support for long-running detection jobs
- Caching for frequently accessed model schemas
- WebSocket support for real-time task updates

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

**Last Updated**: 2025-11-14
**Status**: ✅ Full implementation complete - 13 MCP tools + embedded documentation resources

**Summary**:
- Created 4 model configuration tools (list, schema, validate, generate)
- Created 5 anomaly detection task tools (create, status, list, cancel, limits)
- Created 2 query tools (query_metrics, get_buildinfo) + 1 health check
- Created 1 documentation search tool (search_docs)
- Implemented MCP resources with embedded vmanomaly documentation:
  - 21 markdown files + 30+ images (~5.4 MB)
  - Bleve full-text search with fuzzy matching
  - Markdown chunking (65KB chunks, 8KB overlap)
  - ~50+ documentation resources
- Implemented struct-based MCP schema pattern with `WithInputSchema[T]()` and `NewTypedToolHandler`
- All API types defined in `internal/vmanomaly/types.go`
- Complete API client with 12 methods in `internal/vmanomaly/client.go`
- Tools organized in separate files: models.go, tasks.go, query.go, info.go, config.go, compatibility.go, alerts.go, docs.go
- Resources organized in separate package: internal/resources/
- Full documentation in README.md with parameter details and examples
- Binary builds successfully: `go build -o bin/mcp-vmanomaly ./cmd/mcp-vmanomaly`
- Binary size: ~28 MB (includes embedded documentation)
