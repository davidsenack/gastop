# Architecture

## Overview

gastop is a terminal UI for Gas Town built with Go and [tview](https://github.com/rivo/tview).

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                             │
│                    (entry point, flags)                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                         tui/app.go                          │
│                   (Application orchestrator)                │
│  - Initializes tview.Application                            │
│  - Sets up layout with panels                               │
│  - Handles global key bindings                              │
│  - Manages refresh loop                                     │
└─────────────────────────────────────────────────────────────┘
         │              │              │              │
         ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ tui/convoys  │ │ tui/beads    │ │ tui/polecats │ │ tui/events   │
│   Panel      │ │   Panel      │ │   Panel      │ │   Panel      │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
         │              │              │              │
         └──────────────┴──────────────┴──────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      adapter/adapter.go                     │
│                    (CLI command executor)                   │
│  - Runs gt/bd commands                                      │
│  - Parses JSON output                                       │
│  - Caches results                                           │
│  - Handles errors gracefully                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       model/models.go                       │
│                    (Data structures)                        │
│  - Convoy, Bead, Polecat, Event                             │
│  - Matches JSON schema from gt/bd                           │
└─────────────────────────────────────────────────────────────┘
```

## Package Structure

```
gastop/
├── cmd/gastop/
│   └── main.go           # Entry point, CLI flags
├── internal/
│   ├── adapter/
│   │   ├── adapter.go    # CLI execution and caching
│   │   ├── convoy.go     # Convoy-specific commands
│   │   ├── bead.go       # Bead-specific commands
│   │   ├── polecat.go    # Polecat-specific commands
│   │   └── events.go     # Event stream handling
│   ├── model/
│   │   ├── convoy.go     # Convoy struct
│   │   ├── bead.go       # Bead struct
│   │   ├── polecat.go    # Polecat struct
│   │   └── event.go      # Event struct
│   ├── tui/
│   │   ├── app.go        # Main application
│   │   ├── layout.go     # Panel layout
│   │   ├── convoys.go    # Convoy panel
│   │   ├── beads.go      # Beads panel
│   │   ├── polecats.go   # Polecats panel
│   │   ├── events.go     # Events panel
│   │   ├── statusbar.go  # Top status bar
│   │   ├── help.go       # Help modal
│   │   └── keys.go       # Key bindings
│   ├── config/
│   │   └── config.go     # Configuration loading
│   └── stuck/
│       └── detector.go   # Stuck work detection
├── docs/
│   ├── ARCHITECTURE.md
│   └── DATA_SOURCES.md
├── scripts/
│   └── build.sh
├── tests/
│   └── fixtures/         # Sample JSON outputs for testing
├── go.mod
├── go.sum
└── README.md
```

## Data Flow

### Refresh Loop

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Timer     │────▶│   Adapter   │────▶│   Panels    │
│  (5s tick)  │     │  (fetch)    │     │  (update)   │
└─────────────┘     └─────────────┘     └─────────────┘
                          │
                          ▼
                    ┌─────────────┐
                    │   Cache     │
                    │ (fallback)  │
                    └─────────────┘
```

1. Timer fires every N seconds (configurable)
2. Adapter executes CLI commands in parallel
3. JSON responses parsed into model structs
4. Results cached (for stale fallback)
5. Panels update their display
6. If command fails, use cached data + show stale indicator

### User Input Flow

```
Key Press → App.InputHandler → Panel.HandleKey → Action
                                                    │
                                    ┌───────────────┼───────────────┐
                                    ▼               ▼               ▼
                              Navigation       Commands        Modals
                            (↑↓ Tab Enter)   (r s o p)       (/ f ?)
```

## Component Details

### Adapter Layer

The adapter layer abstracts CLI execution:

```go
type Adapter struct {
    gtPath    string
    bdPath    string
    townRoot  string
    cache     *Cache
    timeout   time.Duration
}

func (a *Adapter) ListConvoys(ctx context.Context) ([]model.Convoy, error)
func (a *Adapter) ListBeads(ctx context.Context, opts BeadListOpts) ([]model.Bead, error)
func (a *Adapter) ListPolecats(ctx context.Context, rig string) ([]model.Polecat, error)
func (a *Adapter) TailEvents(ctx context.Context, n int) ([]model.Event, error)
```

Features:
- Context-based cancellation
- Timeout handling (default 5s per command)
- JSON parsing with unknown field tolerance
- LRU cache with TTL
- Error aggregation

### TUI Panels

Each panel is a tview primitive with:

```go
type Panel interface {
    // Primitive returns the tview component
    Primitive() tview.Primitive

    // Update refreshes data from adapter
    Update(data interface{})

    // HandleKey processes panel-specific keys
    HandleKey(event *tcell.EventKey) bool

    // Selected returns currently selected item
    Selected() interface{}

    // SetFocus handles focus changes
    SetFocus(focused bool)
}
```

### Stuck Detector

Runs on each refresh, checks:

1. **Stale in_progress beads**:
   ```go
   if bead.Status == "in_progress" &&
      time.Since(bead.UpdatedAt) > stuckThreshold {
       bead.Stuck = true
       bead.StuckReason = "No updates for 30+ minutes"
   }
   ```

2. **Orphaned polecats**:
   ```go
   if polecat.State == "working" &&
      polecat.AssignedBead == "" {
       polecat.Stuck = true
       polecat.StuckReason = "Working but no assigned bead"
   }
   ```

3. **Heartbeat timeout** (if available):
   ```go
   if polecat.LastActivity.Before(time.Now().Add(-heartbeatThreshold)) {
       polecat.Stuck = true
       polecat.StuckReason = "No heartbeat"
   }
   ```

## Performance Considerations

### Command Execution

- Run independent commands in parallel (convoys, beads, polecats)
- Use `--limit` flags to bound result sizes
- Cancel in-flight commands on quit

### Rendering

- tview handles efficient terminal updates
- Only redraw panels with changed data
- Use table virtualization for large lists

### Memory

- LRU cache with max size
- Clear old events from stream
- Reuse model structs where possible

## Adding New Panels

1. Create `tui/newpanel.go` implementing Panel interface
2. Add data fetch method to Adapter
3. Register panel in `tui/layout.go`
4. Add key bindings in `tui/keys.go`
5. Update refresh loop in `tui/app.go`

## Configuration

Loaded from (in order):
1. `~/.config/gastop/config.toml`
2. `$XDG_CONFIG_HOME/gastop/config.toml`
3. Command-line flags (override)

```go
type Config struct {
    RefreshInterval     time.Duration
    StuckThresholdMins  int
    LogLines            int
    ShowLogs            bool
    GTPath              string
    BDPath              string
    TownRoot            string
    DefaultFilters      FilterConfig
}
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| gt not in PATH | Show error modal, exit |
| Command timeout | Use cache, show ⚠ stale |
| JSON parse error | Log warning, skip item |
| Permission denied | Show in status bar |
| Network error | N/A (all local) |

## Testing Strategy

1. **Unit tests**: Model parsing, stuck detection logic
2. **Fixture tests**: Parse sample JSON outputs
3. **Integration tests**: Mock adapter with fixtures
4. **Manual tests**: Run against live Gas Town

Test fixtures in `tests/fixtures/`:
- `convoy_list.json`
- `bead_list.json`
- `polecat_list.json`
- `events.jsonl`
