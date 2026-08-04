#!/usr/bin/env bash
# Public installer: downloads the latest precompiled release binary.
# No Go toolchain required — see dev-install.sh if you're building from source.
set -euo pipefail

REPO="moinmalik996/fahh_alert"
BIN_NAME="fahh_alert"
ASSET="fahh_alert-darwin-universal"
INSTALL_DIR="$HOME/.fahh_alert/bin"
ZSHRC="$HOME/.zshrc"

if [ "$(uname -s)" != "Darwin" ]; then
  echo "fahh_alert currently only supports macOS." >&2
  exit 1
fi

echo "Downloading latest ${BIN_NAME} release..."
mkdir -p "$INSTALL_DIR"
curl -fsSL "https://github.com/${REPO}/releases/latest/download/${ASSET}" -o "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

MARKER_START="# >>> fahh_alert hook >>>"
MARKER_END="# <<< fahh_alert hook <<<"

if [ -f "$ZSHRC" ] && grep -qF "$MARKER_START" "$ZSHRC"; then
  echo "Hook already present in ${ZSHRC}, skipping."
else
  echo "Adding hook to ${ZSHRC}..."
  cat >> "$ZSHRC" <<EOF

$MARKER_START
autoload -Uz add-zsh-hook

fahh_alert_precmd() {
    local exit_code=\$?
    "$INSTALL_DIR/$BIN_NAME" "\$exit_code"
}

add-zsh-hook precmd fahh_alert_precmd
$MARKER_END
EOF
fi

echo "Done. Restart your terminal or run: source ${ZSHRC}"
