#!/bin/bash
set -euo pipefail

# Build script for gastop

VERSION="${VERSION:-dev}"
COMMIT="${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')}"
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT}"

echo "Building gastop ${VERSION} (${COMMIT})..."

# Build for current platform
go build -ldflags "${LDFLAGS}" -o gastop ./cmd/gastop

echo "Built: ./gastop"

# Optionally build for multiple platforms
if [[ "${CROSS_COMPILE:-}" == "1" ]]; then
    echo "Cross-compiling..."

    mkdir -p dist

    GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/gastop-linux-amd64 ./cmd/gastop
    GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/gastop-linux-arm64 ./cmd/gastop
    GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/gastop-darwin-amd64 ./cmd/gastop
    GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/gastop-darwin-arm64 ./cmd/gastop

    echo "Cross-compilation complete. Binaries in dist/"
fi

echo "Done!"
