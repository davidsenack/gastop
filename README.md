# gastop

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

An htop-like terminal UI for [Gas Town](https://github.com/anthropics/gas-town) workspaces. Monitor convoys, beads, and polecats in real time.

## Features

- **Convoy Dashboard** - Track batched work with progress bars
- **Bead Browser** - View and filter issues by status, search by ID/title
- **Polecat Monitor** - See worker agents with activity time and hooked work
- **Event Stream** - Real-time activity log
- **Stuck Detection** - Highlights stalled work automatically
- **Vim Navigation** - `j/k/h/l` keys for fast navigation
- **Fast Refresh** - 1-second updates with animated spinners

## Quick Start

### Installation

```bash
git clone https://github.com/davidsenack/gastop.git
cd gastop
go build -o gastop ./cmd/gastop

# Optional: install to PATH
sudo mv gastop /usr/local/bin/
```

### Requirements

- Go 1.21+
- Gas Town CLI tools (`gt` and `bd`) in your PATH

## Using with Your Project

gastop automatically detects your Gas Town workspace. There are several ways to use it:

### Option 1: Run from anywhere (auto-detect)

gastop looks for a Gas Town workspace by:
1. Checking `GT_TOWN_ROOT` environment variable
2. Walking up from current directory looking for `.beads/` or `mayor/` folders
3. Checking common locations (`~/gt/gastop`, `~/gt`, `~/gastop`)

```bash
# If you're inside a Gas Town workspace
cd ~/my-project
gastop

# Or set the environment variable
export GT_TOWN_ROOT=~/my-project
gastop
```

### Option 2: Specify town root explicitly

```bash
gastop -town ~/my-project
```

### Option 3: Create an alias

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
alias gastop='gastop -town ~/my-project'
```

### Option 4: Use a wrapper script

Create a script in your project:

```bash
#!/bin/bash
# my-project/gastop-run.sh
exec gastop -town "$(dirname "$0")"
```

## Keyboard Controls

| Key | Action |
|-----|--------|
| `j/k` | Navigate up/down |
| `h/l` | Switch panels |
| `Tab` | Next panel |
| `Enter` | Select/drill down |
| `/` | Search beads |
| `f` | Filter by status |
| `x` or `d` | Kill polecat / Close bead |
| `r` | Manual refresh |
| `+/-` | Adjust refresh speed |
| `g/G` | Jump to top/bottom |
| `?` | Show help |
| `q` | Quit |

## Layout

```
┌─────────────────────────────────────────────────────────────────────┐
│ gastop │ Town: myproject │ ↻ 1s ◐ │ 14:32:15 │ ? help              │
├───────────────────┬────────────────────────┬────────────────────────┤
│ CONVOYS           │ BEADS                  │ POLECATS               │
│                   │                        │                        │
│ ● Sprint 1        │   ID      Status Title │ ⠹ myproj/furiosa (2m)  │
│   mp-abc [████░░] │ ● mp-123  open   Auth  │   mp-xyz: Add login    │
│   3/5             │ ● mp-124  doing  API   │                        │
│                   │ ✓ mp-125  done   Tests │ ✓ myproj/nux (15m)     │
│                   │                        │   done                 │
├───────────────────┴────────────────────────┴────────────────────────┤
│ EVENTS                                                              │
│ 14:32:01 ● spawn myproj/furiosa                                     │
│ 14:32:02 → sling mp-xyz to furiosa                                  │
│ 14:32:15 ✓ mp-125 completed                                         │
└─────────────────────────────────────────────────────────────────────┘
```

## Status Indicators

| Icon | Meaning |
|------|---------|
| `⠹` (spinner) | Working |
| `✓` | Done/Completed |
| `●` | Open/Active |
| `⚠` | Stuck/Error |
| `○` | Idle |

**Polecat indicators:**
- Red dot `●` after name = session stopped
- Blue ring `◉` = session attached (someone watching)
- `(5m)` = time since last activity

## Configuration

Create `~/.config/gastop/config.toml`:

```toml
# Refresh interval (default: 1s)
refresh_interval = "1s"

# Stuck detection threshold in minutes (default: 30)
stuck_threshold_minutes = 30

# Number of event log lines to show
log_lines = 10

# Show logs panel by default
show_logs = true

[paths]
# Override CLI paths if needed
gt_binary = "gt"
bd_binary = "bd"

# Explicit town root (empty = auto-detect)
town_root = ""

[filters]
# Default status filter for beads
status = ["open", "in_progress"]
show_closed = false
```

## Troubleshooting

### "No active polecats" / Empty panels

1. Make sure you're in a Gas Town workspace or specify `-town`
2. Check that `gt` and `bd` commands work: `gt status`, `bd list`
3. Verify your workspace has a `.beads/` directory

### gastop can't find workspace

```bash
# Check if auto-detection works
gastop -town /path/to/your/project

# Or set environment variable
export GT_TOWN_ROOT=/path/to/your/project
gastop
```

### Slow refresh / timeouts

The default command timeout is 5 seconds. If your `gt`/`bd` commands are slow:

```bash
# Check command speed
time gt polecat list --all
time bd list --limit 10
```

## Development

```bash
# Run directly
go run ./cmd/gastop

# Run tests
go test ./...

# Build
go build -o gastop ./cmd/gastop
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [tview](https://github.com/rivo/tview) - Terminal UI library
- [Gas Town](https://github.com/anthropics/gas-town) - Multi-agent workspace manager
