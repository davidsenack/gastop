#!/bin/bash
# Build a .deb package for gastop
# Usage: ./build-deb.sh [version]

set -e

VERSION="${1:-0.1.0}"
PACKAGE_NAME="gastop"
ARCH=$(dpkg --print-architecture)
BUILD_DIR="$(mktemp -d)"
PACKAGE_DIR="${BUILD_DIR}/${PACKAGE_NAME}_${VERSION}_${ARCH}"

echo "Building ${PACKAGE_NAME} v${VERSION} for ${ARCH}..."

# Create directory structure
mkdir -p "${PACKAGE_DIR}/DEBIAN"
mkdir -p "${PACKAGE_DIR}/usr/bin"
mkdir -p "${PACKAGE_DIR}/usr/share/doc/${PACKAGE_NAME}"
mkdir -p "${PACKAGE_DIR}/usr/share/man/man1"

# Build the binary
cd "$(dirname "$0")/../.."
go build -ldflags "-s -w -X main.version=${VERSION}" -o "${PACKAGE_DIR}/usr/bin/gastop" ./cmd/gastop

# Create control file
cat > "${PACKAGE_DIR}/DEBIAN/control" << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: David Senack <david.senack@gmail.com>
Description: htop-like terminal UI for Gas Town workspaces
 gastop is a terminal-based monitoring tool for Gas Town multi-agent
 workspaces. It provides real-time visualization of convoys, beads,
 and polecats with vim-style keyboard navigation.
Homepage: https://github.com/davidsenack/gastop
EOF

# Copy documentation
cp LICENSE "${PACKAGE_DIR}/usr/share/doc/${PACKAGE_NAME}/copyright" 2>/dev/null || echo "MIT License" > "${PACKAGE_DIR}/usr/share/doc/${PACKAGE_NAME}/copyright"
cp README.md "${PACKAGE_DIR}/usr/share/doc/${PACKAGE_NAME}/"

# Create changelog
cat > "${PACKAGE_DIR}/usr/share/doc/${PACKAGE_NAME}/changelog.Debian" << EOF
${PACKAGE_NAME} (${VERSION}) stable; urgency=low

  * Initial release

 -- David Senack <david.senack@gmail.com>  $(date -R)
EOF
gzip -9 "${PACKAGE_DIR}/usr/share/doc/${PACKAGE_NAME}/changelog.Debian"

# Set permissions
chmod 755 "${PACKAGE_DIR}/usr/bin/gastop"

# Build the package
dpkg-deb --build "${PACKAGE_DIR}"

# Move to current directory
mv "${BUILD_DIR}/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb" .

# Cleanup
rm -rf "${BUILD_DIR}"

echo "Created: ${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
