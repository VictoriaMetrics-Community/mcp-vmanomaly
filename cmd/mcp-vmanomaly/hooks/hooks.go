package hooks

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/VictoriaMetrics/metrics"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func New(ms *metrics.Set) *server.Hooks {
	hooks := &server.Hooks{}

	hooks.AddAfterInitialize(func(_ context.Context, _ any, message *mcp.InitializeRequest, _ *mcp.InitializeResult) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_vmanomaly_initialize_total{client_name="%s",client_version="%s"}`,
			message.Params.ClientInfo.Name,
			message.Params.ClientInfo.Version,
		)).Inc()
	})

	hooks.AddAfterListTools(func(_ context.Context, _ any, _ *mcp.ListToolsRequest, _ *mcp.ListToolsResult) {
		ms.GetOrCreateCounter(`mcp_vmanomaly_list_tools_total`).Inc()
	})

	hooks.AddAfterListResources(func(_ context.Context, _ any, _ *mcp.ListResourcesRequest, _ *mcp.ListResourcesResult) {
		ms.GetOrCreateCounter(`mcp_vmanomaly_list_resources_total`).Inc()
	})

	hooks.AddAfterListPrompts(func(_ context.Context, _ any, _ *mcp.ListPromptsRequest, _ *mcp.ListPromptsResult) {
		ms.GetOrCreateCounter(`mcp_vmanomaly_list_prompts_total`).Inc()
	})

	hooks.AddAfterCallTool(func(_ context.Context, _ any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_vmanomaly_call_tool_total{name="%s",is_error="%t"}`,
			message.Params.Name,
			result.IsError,
		)).Inc()

		if result.IsError {
			slog.Error("Tool call failed", "tool", message.Params.Name)
		} else {
			slog.Info("Tool called", "tool", message.Params.Name, "args", message.Params.Arguments)
			slog.Debug("Tool response", "tool", message.Params.Name, "content", result.Content)
		}
	})

	hooks.AddAfterGetPrompt(func(_ context.Context, _ any, message *mcp.GetPromptRequest, _ *mcp.GetPromptResult) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_vmanomaly_get_prompt_total{name="%s"}`,
			message.Params.Name,
		)).Inc()
	})

	hooks.AddAfterReadResource(func(_ context.Context, _ any, message *mcp.ReadResourceRequest, _ *mcp.ReadResourceResult) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_vmanomaly_read_resource_total{uri="%s"}`,
			message.Params.URI,
		)).Inc()
	})

	hooks.AddOnError(func(_ context.Context, _ any, method mcp.MCPMethod, _ any, err error) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_vmanomaly_error_total{method="%s",error="%s"}`,
			method,
			err,
		)).Inc()

		slog.Error("MCP operation error", "method", method, "error", err)
	})

	return hooks
}
