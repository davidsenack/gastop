# Data Sources

This document describes the data sources gastop uses to display Gas Town state.

## Design Philosophy

1. **CLI-first**: Prefer `gt` and `bd` CLI commands over direct file access
2. **JSON when available**: All key commands support `--json` output
3. **Fallback to file**: Only read files directly if CLI is insufficient
4. **Tolerant parsing**: Handle missing fields, version changes gracefully

---

## CLI Commands (Primary Data Sources)

### Town Status

```bash
gt status --json
```

Returns overall workspace state including:
- Town name and path
- Registered rigs
- Running agents (mayor, deacon)
- Agent states

### Convoy List

```bash
gt convoy list --json
gt convoy list --all --json       # Include closed
gt convoy list --status=closed    # Recently landed
```

Returns:
- Convoy ID (hq-* prefix)
- Title
- Status (open/closed)
- Tracked issues count
- Completion progress

### Convoy Status

```bash
gt convoy status <convoy-id> --json
```

Returns detailed convoy info:
- Tracked issues with status
- Active workers (swarm)
- Completion timestamps

### Bead List

```bash
bd list --json --limit 100
bd list --status=open --json
bd list --status=in_progress --json
bd ready --json                    # Ready work (no blockers)
bd blocked --json                  # Blocked issues
```

Returns issues with:
- `id`: Issue ID with prefix (e.g., `gt-abc`, `bd-xyz`)
- `title`: Issue title
- `status`: open, in_progress, blocked, deferred, closed
- `priority`: 0-4 (0=highest)
- `issue_type`: bug, feature, task, epic, molecule, agent, etc.
- `owner`: Email or name
- `assignee`: Assigned agent/person
- `created_at`, `updated_at`, `closed_at`: Timestamps (RFC3339)
- `labels`: Array of labels
- `description`: Issue description
- Dependencies (via `bd show --json`)

### Bead Show

```bash
bd show <id> --json
```

Returns full issue details including:
- All fields from list
- `blocks`: Issues this blocks
- `blocked_by`: Issues blocking this
- `parent`: Parent epic/molecule
- `children`: Child issues

### Polecat List

```bash
gt polecat list <rig> --json
gt polecat list --all --json       # All rigs
```

Returns active polecats with:
- Name (e.g., "Toast", "Furiosa")
- Rig
- State: working, done, stuck, idle
- Assigned issue ID
- Session status

### Polecat Status

```bash
gt polecat status <rig>/<polecat> --json
```

Returns detailed polecat info:
- Lifecycle state
- Assigned issue
- Session running/attached status
- Creation time
- Last activity time

### Rig List

```bash
gt rig list
```

Returns registered rigs (no JSON yet - parse text output):
- Rig name
- Path
- Prefix
- State (active, parked, docked)

### Logs

```bash
gt log -n 50                       # Last 50 events
gt log --type spawn                # Filter by type
gt log --since 1h                  # Last hour
gt log -f                          # Follow mode
```

Event types: spawn, wake, nudge, handoff, done, crash, kill

No JSON output - parse text. Format:
```
2026-01-22 06:14:21 spawn greenplace/Toast
```

---

## File-Based Data (Secondary Sources)

Use only when CLI is insufficient or for real-time streaming.

### Events File

**Path**: `~/gt/.events.jsonl`

JSONL format, one event per line:
```json
{
  "ts": "2026-01-16T15:30:13Z",
  "source": "gt",
  "type": "session_start",
  "actor": "mayor",
  "payload": {
    "actor_pid": "mayor-18592",
    "cwd": "/Users/davidsenack/gt/mayor",
    "role": "mayor",
    "session_id": "mayor-18592"
  },
  "visibility": "feed"
}
```

Event types:
- `session_start`: Agent session started
- `spawn`: Polecat spawned
- `sling`: Work assigned to agent
- `handoff`: Session handed off
- Various MQ events

**Use case**: Real-time event tail, historical activity

### Issues File

**Path**: `<rig>/.beads/issues.jsonl` (or town-level `.beads/issues.jsonl`)

JSONL format, one issue per line. Same schema as `bd list --json` output.

**Use case**: Direct read if daemon is down, batch processing

### Routes File

**Path**: `~/gt/.beads/routes.jsonl`

Maps issue prefixes to rig paths for routing.

---

## Data Refresh Strategy

### Polling Intervals

| Data Type | Default Interval | Rationale |
|-----------|------------------|-----------|
| Convoys | 5s | User's primary view |
| Beads | 5s | Active work tracking |
| Polecats | 3s | Agent state changes quickly |
| Events | 1s (tail) | Near real-time |
| Town Status | 10s | Rarely changes |

### Caching

- Cache last successful response per command
- On error, return stale data with "stale" indicator
- Clear cache on manual refresh (r key)

### Diffing

- Compare new data with cached data
- Only update UI elements that changed
- Preserve scroll position when possible

---

## Stuck Detection Data

To detect stuck work, combine:

1. **Bead status + updated_at**:
   - `status: in_progress` AND `updated_at` > X minutes ago → stuck
   - Query: `bd list --status=in_progress --json`

2. **Polecat state + last activity**:
   - Polecat assigned but no heartbeat/activity
   - Query: `gt polecat status <rig>/<polecat> --json`

3. **Orphaned work**:
   - `gt orphans` - Find lost polecat work

Default stuck threshold: 30 minutes (configurable)

---

## Command Error Handling

| Error | Action |
|-------|--------|
| Command not found | Show error toast, disable feature |
| Timeout | Use cached data, mark stale |
| Parse error | Log warning, skip record |
| Permission denied | Show error, suggest fix |

---

## Version Compatibility

- Require `gt` and `bd` in PATH
- Check `gt version` and `bd version` on startup
- Warn if version mismatch detected
- Gracefully handle missing JSON fields (schema evolution)
