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

## What You See

```
┌──────────────────────────────────────────────────────────────┐
│ gastop │ Town: myproject │ ↻ 1s │ 14:32:15 │ ? help         │
├─────────────────┬──────────────────┬─────────────────────────┤
│ CONVOYS         │ BEADS            │ POLECATS                │
│ ● Sprint 1      │ ● mp-123  Auth   │ ⠹ furiosa (2m)          │
│   [████░░] 3/5  │ ● mp-124  API    │   mp-xyz: Add login     │
│                 │ ✓ mp-125  Tests  │ ✓ nux (15m) done        │
├─────────────────┴──────────────────┴─────────────────────────┤
│ EVENTS                                                       │
│ 14:32:01 ● spawn furiosa                                     │
│ 14:32:15 ✓ mp-125 completed                                  │
└──────────────────────────────────────────────────────────────┘
```

- **Convoys** - Batched work with progress bars
- **Beads** - Issues/tasks with status
- **Polecats** - Worker agents with activity time
- **Events** - Real-time activity log

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

**Homebrew:**
```bash
brew install davidsenack/tap/gastop
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
