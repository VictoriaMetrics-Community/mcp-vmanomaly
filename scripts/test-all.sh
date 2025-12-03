#!/bin/bash
set -e
set -o pipefail

# Run tests
go test -v -race -coverprofile=coverage.out ./...
