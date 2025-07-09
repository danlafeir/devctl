#!/bin/sh
set -e

REPO=danlafeir/devctl
BINARY=devctl
INSTALL_DIR=/usr/local/bin

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux) OS=linux ;;
  darwin) OS=darwin ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect ARCH
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version tag from GitHub API
LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name"' | cut -d '"' -f4)
if [ -z "$LATEST" ]; then
  echo "Could not determine latest release version." >&2
  exit 1
fi

FILENAME="$BINARY-$OS-$ARCH"
URL="https://github.com/$REPO/releases/download/$LATEST/$FILENAME"

TMP=$(mktemp)
echo "Downloading $URL ..."
curl -sSLfL "$URL" -o "$TMP"
chmod +x "$TMP"

# Move to install dir (may require sudo)
echo "Installing to $INSTALL_DIR/$BINARY ..."
sudo mv "$TMP" "$INSTALL_DIR/$BINARY"

if command -v $BINARY >/dev/null 2>&1; then
  echo "$BINARY installed successfully!"
  $BINARY --help
else
  echo "Install failed: $BINARY not found in PATH." >&2
  exit 1
fi 