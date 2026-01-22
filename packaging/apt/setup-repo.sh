#!/bin/bash
# Script to set up the apt repository structure
# This is run by GitHub Actions after building .deb packages

set -e

REPO_DIR="${1:-apt-repo}"
VERSION="${2:-0.1.0}"

mkdir -p "$REPO_DIR/pool/main"
mkdir -p "$REPO_DIR/dists/stable/main/binary-amd64"
mkdir -p "$REPO_DIR/dists/stable/main/binary-arm64"

# Move deb files to pool
mv *.deb "$REPO_DIR/pool/main/" 2>/dev/null || true

# Generate Packages files
cd "$REPO_DIR"

apt-ftparchive packages pool/main > dists/stable/main/binary-amd64/Packages
gzip -k dists/stable/main/binary-amd64/Packages

apt-ftparchive packages pool/main > dists/stable/main/binary-arm64/Packages
gzip -k dists/stable/main/binary-arm64/Packages

# Generate Release file
cat > dists/stable/Release << EOF
Origin: gastop
Label: gastop
Suite: stable
Codename: stable
Architectures: amd64 arm64
Components: main
Description: gastop - Terminal dashboard for Gas Town
EOF

apt-ftparchive release dists/stable >> dists/stable/Release

echo "Apt repository created in $REPO_DIR"
