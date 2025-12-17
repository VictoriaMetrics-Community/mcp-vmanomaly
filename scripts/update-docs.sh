#!/bin/bash
set -e
set -o pipefail

# Update vmanomaly documentation from VictoriaMetrics docs repository

# Remove existing docs
rm -rf internal/resources/docs/anomaly-detection

# Clone VictoriaMetrics docs repository with sparse checkout
git clone --no-checkout --depth=1 https://github.com/VictoriaMetrics/vmdocs.git /tmp/vmdocs-temp
cd /tmp/vmdocs-temp

# Setup sparse checkout to get only anomaly-detection folder
git sparse-checkout init --cone
git sparse-checkout set content/docs/anomaly-detection
git checkout main

# Copy anomaly-detection docs to our resources
mkdir -p $(pwd -P | sed 's|/tmp/vmdocs-temp||')/internal/resources/docs
cp -r content/docs/anomaly-detection $(pwd -P | sed 's|/tmp/vmdocs-temp||')/internal/resources/docs/

# Cleanup
cd -
rm -rf /tmp/vmdocs-temp

echo "✅ Documentation updated successfully!"
echo "📁 Location: internal/resources/docs/anomaly-detection"
