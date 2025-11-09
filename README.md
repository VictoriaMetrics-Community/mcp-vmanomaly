# MCP Server for vmanomaly

Model Context Protocol (MCP) server for [vmanomaly](https://docs.victoriametrics.com/anomaly-detection/) - VictoriaMetrics' anomaly detection application.

This MCP server enables AI assistants like Claude to interact with vmanomaly's REST API for anomaly detection, model management, and observability insights.

## Features

- **Health Monitoring**: Check vmanomaly server health and status
- **Model Management**: List and manage anomaly detection models (coming soon)
- **Anomaly Detection**: Trigger and retrieve anomaly detection results (coming soon)
- **Metrics & Logs**: Query VictoriaMetrics and VictoriaLogs for anomaly insights (coming soon)

## Prerequisites

- Go 1.24 or higher
- A running vmanomaly instance
- Access to vmanomaly REST API endpoint

## Installation

### From Source

```bash
# Clone the repository
git clone <repository-url>
cd mcp-vmanomaly

# Install dependencies
make install

# Build the binary
make build

# The binary will be available at ./bin/mcp-vmanomaly
```

### Using Go Install

```bash
go install github.com/VictoriaMetrics/mcp-vmanomaly/cmd/mcp-vmanomaly@latest
```

## Configuration

The MCP server is configured using environment variables:

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `VMANOMALY_ENDPOINT` | Yes | vmanomaly server endpoint URL | `http://localhost:8490` |
| `VMANOMALY_API_TOKEN` | No | API token for authentication (if required) | `your-api-token` |

## Usage

### Running the Server

```bash
# Set required environment variables
export VMANOMALY_ENDPOINT=http://localhost:8490
export VMANOMALY_API_TOKEN=your-token  # Optional

# Run the server
./bin/mcp-vmanomaly
```

Or using make:

```bash
VMANOMALY_ENDPOINT=http://localhost:8490 make run
```

### Integration with Claude Desktop

Add the following to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "vmanomaly": {
      "command": "/path/to/bin/mcp-vmanomaly",
      "env": {
        "VMANOMALY_ENDPOINT": "http://localhost:8490",
        "VMANOMALY_API_TOKEN": "your-token-if-needed"
      }
    }
  }
}
```

### Integration with Other MCP Clients

#### Cursor

Add to `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "vmanomaly": {
      "command": "/path/to/bin/mcp-vmanomaly",
      "env": {
        "VMANOMALY_ENDPOINT": "http://localhost:8490"
      }
    }
  }
}
```

#### VS Code (with Cline extension)

Add to VS Code settings or Cline configuration:

```json
{
  "mcp.servers": {
    "vmanomaly": {
      "command": "/path/to/bin/mcp-vmanomaly",
      "env": {
        "VMANOMALY_ENDPOINT": "http://localhost:8490"
      }
    }
  }
}
```

## Available Tools

### `health_check`

Check the health status of the vmanomaly server.

**Parameters**:
- `endpoint` (optional): Override the configured vmanomaly endpoint

**Example**:
```
Check the health of my vmanomaly server
```

### Coming Soon

- `list_models`: List available anomaly detection models
- `run_detection`: Trigger anomaly detection
- `get_anomalies`: Retrieve detected anomalies
- `get_model_metrics`: Get model performance metrics
- `query_logs`: Query VictoriaLogs for anomaly-related data

## Development

### Project Structure

```
mcp-vmanomaly/
├── cmd/
│   └── mcp-vmanomaly/      # Main application entry point
│       └── main.go
├── internal/
│   ├── vmanomaly/          # vmanomaly API client
│   │   └── client.go
│   └── tools/              # MCP tool definitions
│       └── tools.go
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── CLAUDE.md               # Claude memory file
```

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linters
make lint
```

### Adding New Tools

1. Add the API method to `internal/vmanomaly/client.go`
2. Create a tool definition in `internal/tools/tools.go`
3. Register the tool in the `RegisterTools` function
4. Update documentation

Example:

```go
// In internal/vmanomaly/client.go
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
    // Implementation
}

// In internal/tools/tools.go
func RegisterTools(s *server.MCPServer, client *vmanomaly.Client) {
    // ... existing tools ...

    listModelsTool := mcp.NewTool("list_models",
        mcp.WithDescription("List available anomaly detection models"),
    )
    s.AddTool(listModelsTool, listModelsHandler(client))
}
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

[Specify your license here]

## Related Projects

- [vmanomaly](https://docs.victoriametrics.com/anomaly-detection/) - VictoriaMetrics anomaly detection
- [VictoriaMetrics](https://victoriametrics.com/) - Time series database
- [mcp-victoriametrics](https://github.com/VictoriaMetrics-Community/mcp-victoriametrics) - MCP server for VictoriaMetrics
- [Model Context Protocol](https://modelcontextprotocol.io/) - MCP specification

## Support

For vmanomaly-specific questions, see the [vmanomaly documentation](https://docs.victoriametrics.com/anomaly-detection/).

For MCP server issues, please open an issue in this repository.
