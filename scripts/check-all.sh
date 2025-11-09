#!/bin/bash
set -e
set -o pipefail

# Check licenses
which wwhrd || go install github.com/frapposelli/wwhrd@latest
wwhrd check -f .wwhrd.yml

# Check for vulnerabilities
which govulncheck || go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
