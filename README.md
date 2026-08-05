# shellsound

A tiny terminal companion that plays a sound after each command:
- failure sound for non-zero exit codes
- success sound for zero exit codes

```
$ python45.43
zsh: command not found: python45.43
🔊 shellsound.mp3

$ ls
🔊 Glass.aiff (system sound)
```

## Platform

- macOS
- zsh

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/moinmalik996/shellsound/main/scripts/install.sh | sh
```

Or install from a local clone:

```bash
git clone https://github.com/moinmalik996/shellsound.git
cd shellsound
./scripts/install.sh
```

Reload your shell and test:

```bash
exec zsh
python45.43     # command not found → plays the failure sound
ls              # succeeds → plays the success sound
```

## Uninstall

From any installed zsh session:

```bash
shellsound uninstall
```

Then reload:

```bash
exec zsh
```

## Configuration

On first run, `shellsound` creates `~/.shellsound/config.json` with these
defaults:

```json
{
  "failureSound": "/Users/you/.shellsound/sounds/shellsound.mp3",
  "successSound": "/System/Library/Sounds/Glass.aiff",
  "playOnFailure": true,
  "playOnSuccess": true
}
```

| Field | Type | Meaning |
|---|---|---|
| `failureSound` | string | Path to the sound file played when a command exits non-zero |
| `successSound` | string | Path to the sound file played when a command exits 0 |
| `playOnFailure` | bool | Whether to play a sound on failure at all |
| `playOnSuccess` | bool | Whether to play a sound on success at all |

Any `afplay`-supported format works (`.mp3`, `.wav`, `.aiff`, `.m4a`, ...).

## Development

Run tests:

```bash
go test ./...
```

Install from source:

```bash
./scripts/dev-install.sh
```

## Roadmap

- [ ] Capture and print the failed command text (via zsh's `preexec` hook)
- [ ] Bash / Fish / PowerShell support
- [ ] Duration-based guard for `playOnSuccess` (skip trivial instant commands)
