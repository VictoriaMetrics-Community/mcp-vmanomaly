#!/bin/bash
set -e
set -o pipefail

# Build the binary
mkdir -p bin
go build -o bin/mcp-vmanomaly ./cmd/mcp-vmanomaly
echo "Build complete: bin/mcp-vmanomaly"
