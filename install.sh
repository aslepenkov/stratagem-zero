#!/bin/sh
set -e

REPO="aslepenkov/stratagem-zero"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux) OS="linux" ;;
  darwin) OS="darwin" ;;
  *)
    echo "Error: Unsupported operating system '$OS'." >&2
    exit 1
    ;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported architecture '$ARCH'." >&2
    exit 1
    ;;
esac

# Helper for download
download() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$1"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "$1"
  else
    echo "Error: Neither curl nor wget is installed." >&2
    exit 1
  fi
}

download_file() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$1" -o "$2"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$2" "$1"
  else
    echo "Error: Neither curl nor wget is installed." >&2
    exit 1
  fi
}

echo "Detecting latest release for $REPO..."

# Get latest release tag
TAG=""
if command -v curl >/dev/null 2>&1; then
  TAG="$(curl -fsSL -I "https://github.com/${REPO}/releases/latest" 2>/dev/null | grep -i "^location:" | sed -E 's/.*tag\/([^/?#]+).*/\1/' | tr -d '\r\n')"
fi

if [ -z "$TAG" ]; then
  LATEST_JSON="$(download "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null || true)"
  TAG="$(echo "$LATEST_JSON" | grep '"tag_name":' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
fi

if [ -z "$TAG" ]; then
  TAG="v0.1.0"
fi

echo "Installing stratagem-zero $TAG ($OS/$ARCH)..."

ARCHIVE_NAME="stratagem-zero_${TAG}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${ARCHIVE_NAME}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading from $DOWNLOAD_URL..."
if ! download_file "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE_NAME"; then
  echo "Error: Failed to download $DOWNLOAD_URL" >&2
  exit 1
fi

echo "Extracting archive..."
tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"

if [ ! -f "$TMP_DIR/stratagem-zero" ]; then
  echo "Error: Executable 'stratagem-zero' not found in archive." >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR"
mv "$TMP_DIR/stratagem-zero" "$INSTALL_DIR/stratagem-zero"
chmod +x "$INSTALL_DIR/stratagem-zero"

echo "Successfully installed stratagem-zero to $INSTALL_DIR/stratagem-zero"

case ":$PATH:" in
  *:"$INSTALL_DIR":*) ;;
  *)
    echo ""
    echo "Note: '$INSTALL_DIR' is not in your PATH."
    echo "Add it to your PATH by adding this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    ;;
esac
