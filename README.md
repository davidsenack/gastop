# gastop

Terminal dashboard for [Gas Town](https://github.com/anthropics/gas-town) workspaces. Like htop, but for your multi-agent projects.

![Screenshot](assets/screenshots/gastop.png)

## Quick Start

```bash
# Install
brew install davidsenack/tap/gastop

# Run (auto-detects your workspace)
gastop

# Or specify a workspace
gastop --town ~/my-project
```

## Keyboard

| Key | Action |
|-----|--------|
| `j/k` | Up/down |
| `h/l` | Switch panels |
| `/` | Search |
| `f` | Filter |
| `x` | Kill/close |
| `?` | Help |
| `q` | Quit |

## Install

**Homebrew (macOS/Linux):**
```bash
brew install davidsenack/tap/gastop
```

**Arch Linux (AUR):**
```bash
yay -S gastop
```

**Debian/Ubuntu:**
```bash
curl -LO https://github.com/davidsenack/gastop/releases/latest/download/gastop_amd64.deb
sudo dpkg -i gastop_amd64.deb
```

**From source:**
```bash
go install github.com/davidsenack/gastop/cmd/gastop@latest
```

**Binary:** Download from [releases](https://github.com/davidsenack/gastop/releases)

## Options

```
-t, --town    Workspace directory (auto-detects if not set)
-r, --rig     Focus on specific rig
-j, --json    JSON output for scripting
-V, --version Show version
```

## Requirements

- [Gas Town](https://github.com/anthropics/gas-town) CLI tools (`gt` and `bd`)

## License

MIT
