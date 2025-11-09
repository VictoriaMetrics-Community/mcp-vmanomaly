#!/bin/bash
set -e
set -o pipefail

# Run golangci-lint
which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run --config .golangci.yml
