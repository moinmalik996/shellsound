# fahh_alert

A tiny Go tool that plays a sound after every command you run in your terminal —
one sound when a command fails, another when it succeeds. No wrapper commands,
no changing how you work: install it once and it silently watches every
command in every new terminal tab.

```
$ python45.43
zsh: command not found: python45.43
🔊 fahhh.mp3

$ ls
🔊 Glass.aiff (system sound)
```

## How it works

This is **not** a program that "watches" your terminal from the outside — no
OS provides an API for that. Instead it hooks into zsh's own lifecycle:

```
You run a command
        │
        ▼
zsh's `precmd` hook fires after the command finishes
        │
        ▼
The hook calls the fahh_alert binary with the command's exit code
        │
        ▼
fahh_alert reads ~/.fahh_alert/config.json
        │
        ▼
exit code == 0  → play the configured "success" sound (if enabled)
exit code != 0  → play the configured "failure" sound (if enabled)
        │
        ▼
Sound plays via macOS's `afplay`
```

Zsh calls the binary through a hook function added to `~/.zshrc`
(see [Installation](#installation) below):

```sh
fahh_alert_precmd() {
    local exit_code=$?
    "$HOME/.fahh_alert/bin/fahh_alert" "$exit_code"
}
add-zsh-hook precmd fahh_alert_precmd
```

Detection is purely by **exit code** — not by scanning output text — so it
works identically for every command: shell builtins, Go/Python/Node programs,
`kubectl`, `docker`, anything.

## Project layout

```
fahh_alert/
  cmd/fahh_alert/main.go        entry point: reads exit code, decides which sound to play
  internal/audio/player.go      plays a sound file via afplay
  internal/config/config.go     loads/creates ~/.fahh_alert/config.json
  internal/assets/assets.go     embeds sounds/fahhh.mp3 directly into the binary
  internal/assets/sounds/       the default failure sound, versioned in the repo
  scripts/install.sh            public installer: downloads the release binary, wires the zsh hook
  scripts/dev-install.sh        dev installer: builds from source instead of downloading
  .github/workflows/release.yml builds a universal binary and publishes a GitHub Release on every version tag
```

## Requirements

- **macOS** (uses `afplay`, which ships with macOS)
- **zsh** (the default shell on macOS since Catalina)
- **Nothing else.** The release binary is fully self-contained — see
  [Do I need Go installed?](#do-i-need-go-installed) below.

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/moinmalik996/fahh_alert/main/scripts/install.sh | sh
```

Or clone the repo and run it locally:

```bash
git clone https://github.com/moinmalik996/fahh_alert.git
cd fahh_alert
./scripts/install.sh
```

The script:

1. Downloads the latest precompiled universal binary (arm64 + Intel) from
   [GitHub Releases](https://github.com/moinmalik996/fahh_alert/releases).
2. Installs it to `~/.fahh_alert/bin/fahh_alert` (no `sudo` needed — it
   lives entirely under your home directory).
3. Appends a marked block to `~/.zshrc` that registers the `precmd` hook.

Open a **new** terminal tab/window (or run `source ~/.zshrc` in an existing
one — `.zshrc` is only read when a shell starts, so already-open tabs won't
pick up the change automatically). Then try:

```bash
python45.43     # command not found → plays the failure sound
ls               # succeeds → plays the success sound
```

### Uninstalling

Delete the block between these two lines in `~/.zshrc`:

```
# >>> fahh_alert hook >>>
...
# <<< fahh_alert hook <<<
```

Then optionally remove `~/.fahh_alert/` (binary + config).

## Configuration

On first run, `fahh_alert` creates `~/.fahh_alert/config.json` with these
defaults:

```json
{
  "failureSound": "/Users/you/.fahh_alert/sounds/fahhh.mp3",
  "successSound": "/System/Library/Sounds/Glass.aiff",
  "playOnFailure": true,
  "playOnSuccess": true
}
```

`failureSound` points at a copy of `fahhh.mp3` extracted from inside the
binary itself on first run — the binary carries this file embedded, so it
works the same on any Mac, with nothing to download separately.

| Field | Type | Meaning |
|---|---|---|
| `failureSound` | string | Path to the sound file played when a command exits non-zero |
| `successSound` | string | Path to the sound file played when a command exits 0 |
| `playOnFailure` | bool | Whether to play a sound on failure at all |
| `playOnSuccess` | bool | Whether to play a sound on success at all |

Any audio format `afplay` supports works (`.mp3`, `.wav`, `.aiff`, `.m4a`,
...). macOS ships a set of built-in sounds under `/System/Library/Sounds/`
(e.g. `Basso.aiff`, `Glass.aiff`, `Hero.aiff`) you can use for either field.

Edit the file and changes take effect on the very next command — no need to
rebuild or restart your shell.

> **Heads up:** with `playOnSuccess` enabled, *every* successful command
> plays a sound, including trivial ones like `ls` or `cd`. If that gets
> noisy, set `"playOnSuccess": false`.

## Do I need Go installed?

**No — not to install or run it.** `scripts/install.sh` downloads a
precompiled universal binary (works on both Apple Silicon and Intel Macs)
from GitHub Releases; nothing gets compiled on your machine. The binary has
no runtime dependency on Go, and its default sound is embedded inside it
(via Go's `embed` package at build time), so there's no separate asset file
to fetch either.

Go is only needed if you want to **build the binary yourself** — either to
contribute, or to test local changes before they're released (see
`scripts/dev-install.sh` and [Development](#development) below). Releases
are built automatically by
[`.github/workflows/release.yml`](.github/workflows/release.yml) whenever a
`v*` tag is pushed, so end users never need Go installed at all.

## Development

Run the test suite:

```bash
go test ./...
```

To test local changes before they're released, build and install from
source instead of downloading:

```bash
./scripts/dev-install.sh
```

This builds `fahh_alert` from source and installs it to the same
`~/.fahh_alert/bin/fahh_alert` location the public installer uses, so you
can iterate without waiting for a tagged release.

### Cutting a release

```bash
git tag v0.1.0
git push origin v0.1.0
```

Pushing a `v*` tag triggers the GitHub Actions workflow, which builds the
universal binary and publishes it as a GitHub Release automatically — no
manual build/upload steps.

## Roadmap

- [ ] Capture and print the failed command text (via zsh's `preexec` hook)
- [ ] Bash / Fish / PowerShell support
- [ ] Duration-based guard for `playOnSuccess` (skip trivial instant commands)
