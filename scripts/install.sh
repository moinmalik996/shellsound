#!/usr/bin/env bash
# Public installer: downloads the latest precompiled release binary.
# No Go toolchain required — see dev-install.sh if you're building from source.
set -euo pipefail

REPO="moinabbas/shellsound"
BIN_NAME="shellsound"
ASSET="shellsound-darwin-universal"
INSTALL_DIR="$HOME/.shellsound/bin"
ZSHRC="$HOME/.zshrc"

if [ "$(uname -s)" != "Darwin" ]; then
  echo "shellsound currently only supports macOS." >&2
  exit 1
fi

echo "Downloading latest ${BIN_NAME} release..."
mkdir -p "$INSTALL_DIR"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

if curl -fsSL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BIN_NAME"; then
  chmod +x "$INSTALL_DIR/$BIN_NAME"
else
  echo "No release asset found at ${DOWNLOAD_URL}." >&2
  echo "Falling back to local source build (requires Go)." >&2

  if ! command -v go >/dev/null 2>&1; then
    echo "Go is not installed. Install Go or publish a release asset first." >&2
    exit 1
  fi

  REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
  (cd "$REPO_DIR" && go build -o "$INSTALL_DIR/$BIN_NAME" ./cmd/shellsound)
  chmod +x "$INSTALL_DIR/$BIN_NAME"
fi

MARKER_START="# >>> shellsound hook >>>"
MARKER_END="# <<< shellsound hook <<<"

if [ -f "$ZSHRC" ] && grep -qF "$MARKER_START" "$ZSHRC"; then
  echo "Hook already present in ${ZSHRC}, skipping."
else
  echo "Adding hook to ${ZSHRC}..."
  cat >> "$ZSHRC" <<EOF

$MARKER_START
autoload -Uz add-zsh-hook

shellsound_precmd() {
    local exit_code=\$?
    "$INSTALL_DIR/$BIN_NAME" "\$exit_code"
}

add-zsh-hook precmd shellsound_precmd
$MARKER_END
EOF
fi

echo "Done. Restart your terminal or run: source ${ZSHRC}"
