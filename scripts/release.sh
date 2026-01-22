#!/bin/bash
# Build release binaries for multiple platforms
# Usage: ./scripts/release.sh [version]

set -e

VERSION="${1:-0.1.0}"
PROJECT_NAME="gastop"
BUILD_DIR="dist"
LDFLAGS="-s -w -X main.version=${VERSION}"

echo "Building ${PROJECT_NAME} v${VERSION}..."

# Clean and create build directory
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

# Platforms to build for
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

cd "$(dirname "$0")/.."

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OUTPUT_NAME="${PROJECT_NAME}-${VERSION}-${GOOS}-${GOARCH}"

    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi

    echo "Building ${GOOS}/${GOARCH}..."

    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "${LDFLAGS}" \
        -o "${BUILD_DIR}/${OUTPUT_NAME}" \
        ./cmd/gastop

    # Create archive
    if [ "$GOOS" = "windows" ]; then
        (cd "${BUILD_DIR}" && zip "${PROJECT_NAME}-${VERSION}-${GOOS}-${GOARCH}.zip" "${OUTPUT_NAME}")
        rm "${BUILD_DIR}/${OUTPUT_NAME}"
    else
        (cd "${BUILD_DIR}" && tar -czf "${PROJECT_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz" "${OUTPUT_NAME}")
        rm "${BUILD_DIR}/${OUTPUT_NAME}"
    fi
done

# Generate checksums
echo "Generating checksums..."
(cd "${BUILD_DIR}" && shasum -a 256 *.tar.gz *.zip > checksums.txt)

echo ""
echo "Release artifacts in ${BUILD_DIR}/:"
ls -la "${BUILD_DIR}"

echo ""
echo "Checksums:"
cat "${BUILD_DIR}/checksums.txt"

echo ""
echo "Done! Upload these files to GitHub releases."
