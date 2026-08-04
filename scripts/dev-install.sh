#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_NAME="fahh_alert"
INSTALL_DIR="$HOME/.fahh_alert/bin"
ZSHRC="$HOME/.zshrc"

echo "Building ${BIN_NAME}..."
(cd "$REPO_DIR" && go build -o "$BIN_NAME" ./cmd/fahh_alert)

echo "Installing to ${INSTALL_DIR}..."
mkdir -p "$INSTALL_DIR"
mv "$REPO_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
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
