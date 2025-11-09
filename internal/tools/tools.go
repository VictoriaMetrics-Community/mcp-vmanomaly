package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"mcp-vmanomaly/internal/vmanomaly"
)

// RegisterTools registers all MCP tools with the server
func RegisterTools(s *server.MCPServer, client *vmanomaly.Client) {
	// Register health check tool
	healthTool := mcp.NewTool("health_check",
		mcp.WithDescription("Check the health status of the vmanomaly server"),
	)

	s.AddTool(healthTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Call vmanomaly health endpoint
		health, err := client.GetHealth(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Health check failed: %v", err)), nil
		}

		// Format response
		responseJSON, err := json.MarshalIndent(health, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(responseJSON)), nil
	})

	// TODO: Add more tools here based on vmanomaly API capabilities:
	// - list_models: List available anomaly detection models
	// - run_detection: Trigger anomaly detection on specified data
	// - get_anomalies: Retrieve detected anomalies
	// - get_model_metrics: Get metrics for a specific model
	// - query_logs: Query VictoriaLogs for anomaly-related logs
}

// Example of how to add more tools:
//
// listModelsTool := mcp.NewTool("list_models",
//     mcp.WithDescription("List available anomaly detection models"),
// )
//
// s.AddTool(listModelsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
//     models, err := client.ListModels(ctx)
//     if err != nil {
//         return mcp.NewToolResultError(fmt.Sprintf("Failed to list models: %v", err)), nil
//     }
//
//     responseJSON, err := json.MarshalIndent(models, "", "  ")
//     if err != nil {
//         return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
//     }
//
//     return mcp.NewToolResultText(string(responseJSON)), nil
// })
