# gastop

> **Warning:** This project is 100% vibe-coded and should not be trusted at all. But have fun!

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
# Add the repository
echo "deb [trusted=yes] https://davidsenack.github.io/gastop stable main" | sudo tee /etc/apt/sources.list.d/gastop.list

# Install
sudo apt update && sudo apt install gastop
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

## Contributing

Contributions welcome! Feel free to open issues or submit PRs. This is a fun project - don't overthink it.

## License

MIT
