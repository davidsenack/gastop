# gastop

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

An htop-like terminal UI for [Gas Town](https://github.com/steveyegge/gastown) - visualize convoys, beads, and polecats in real time.

![gastop screenshot](docs/screenshot.png)

## Features

- **Convoy Dashboard** - Track batched work progress across rigs
- **Bead Browser** - View and filter issues by status, priority, type
- **Polecat Monitor** - See running agents and their assigned work
- **Event Stream** - Real-time activity log with filtering
- **Stuck Detector** - Highlights work that hasn't progressed
- **Keyboard-Driven** - htop-style navigation and actions
- **Fast Refresh** - 1-second updates with incremental rendering
- **Resilient** - Graceful degradation with cached data on errors

## Installation

### From Source

```bash
git clone https://github.com/davidsenack/gastop.git
cd gastop
go build -o gastop ./cmd/gastop

# Or install to $GOPATH/bin
go install ./cmd/gastop
```

### Requirements

- Go 1.21+
- [Gas Town](https://github.com/steveyegge/gastown) (`gt` and `bd` CLI tools in PATH)

## Usage

```bash
# Launch TUI
gastop

# Specify town root
gastop --town ~/gt

# Start with specific rig focus
gastop --rig gastown

# JSON output mode (for scripting)
gastop --json
```

## Keyboard Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate lists |
| `Tab` | Switch panes |
| `Enter` | Drill down (convoy → beads → detail) |
| `/` | Search |
| `f` | Filter by status |
| `r` | Manual refresh |
| `t` | Toggle auto-refresh |
| `l` | Toggle logs panel |
| `o` | Open selected in $EDITOR |
| `s` | Sling selected bead |
| `p` | Pause/resume convoy |
| `?` | Show help |
| `q` | Quit |

## Layout

```
┌─────────────────────────────────────────────────────────────────┐
│ gastop | Town: gastop | Rig: gastown | ↻ 5s | ● Connected      │
├──────────────────┬───────────────────────┬──────────────────────┤
│ CONVOYS          │ BEADS                 │ POLECATS             │
│                  │                       │                      │
│ ▸ gastop MVP    │ ID      Status  Title │ Toast    working     │
│   hq-abc [3/5]   │ gt-123  ●doing  Recon │ Furiosa  idle        │
│                  │ gt-124  ○open   Model │ Nux      stuck ⚠     │
│   another-conv   │ gt-125  ✓done   Parse │                      │
│   hq-xyz [2/2]   │                       │                      │
│                  │                       │                      │
├──────────────────┴───────────────────────┴──────────────────────┤
│ EVENTS                                                          │
│ 06:14:21 spawn greenplace/Toast                                 │
│ 06:14:22 sling gt-123 → Toast                                   │
│ 06:15:01 ✓ gt-125 completed                                     │
└─────────────────────────────────────────────────────────────────┘
```

## Configuration

Create `~/.config/gastop/config.toml`:

```toml
# Refresh interval in seconds
refresh_interval = 5

# Stuck detection threshold in minutes
stuck_threshold_minutes = 30

# Default filters
[filters]
status = ["open", "in_progress"]
show_closed = false

# Log panel
[logs]
visible = true
lines = 10

# Paths
[paths]
gt_binary = "gt"
bd_binary = "bd"
town_root = ""  # Auto-detect if empty
```

## Architecture

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for design details.

## Development

```bash
# Run with live reload
go run ./cmd/gastop

# Run tests
go test ./...

# Build release
./scripts/build.sh
```

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [tview](https://github.com/rivo/tview) - Terminal UI library
- [Gas Town](https://github.com/steveyegge/gastown) - Multi-agent workspace manager
