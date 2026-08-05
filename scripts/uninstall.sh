#!/usr/bin/env bash
set -euo pipefail

ZSHRC="$HOME/.zshrc"
SHELLSOUND_DIR="$HOME/.shellsound"

SHELLSOUND_START="# >>> shellsound hook >>>"
SHELLSOUND_END="# <<< shellsound hook <<<"

info() {
  echo "[shellsound-uninstall] $*"
}

remove_hook_block() {
  if [ ! -f "$ZSHRC" ]; then
    info "No ${ZSHRC} found. Skipping hook cleanup."
    return
  fi

  if ! grep -qF "$SHELLSOUND_START" "$ZSHRC"; then
    info "No shellsound hook block found in ${ZSHRC}."
    return
  fi

  local backup
  backup="${ZSHRC}.backup.$(date +%Y%m%d-%H%M%S)"
  cp "$ZSHRC" "$backup"
  info "Backed up ${ZSHRC} to ${backup}"

  perl -0777 -i '' -pe "s/\\n?\Q${SHELLSOUND_START}\E.*?\Q${SHELLSOUND_END}\E\\n?//sg" "$ZSHRC"
  info "Removed shellsound hook block from ${ZSHRC}"
}

remove_install_dir() {
  if [ -d "$SHELLSOUND_DIR" ]; then
    rm -rf "$SHELLSOUND_DIR"
    info "Removed install directory: ${SHELLSOUND_DIR}"
  else
    info "Install directory not found: ${SHELLSOUND_DIR}"
  fi
}

remove_hook_block
remove_install_dir

info "Uninstall complete."
info "Reload your shell with: exec zsh"