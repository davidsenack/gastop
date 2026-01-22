package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

// TailEvents returns the last N events from the activity log.
func (a *Adapter) TailEvents(ctx context.Context, n int) ([]model.Event, error) {
	// Try to read directly from the events file for speed
	eventsPath := filepath.Join(a.townRoot, ".events.jsonl")
	if a.townRoot == "" {
		// Try to find the town root
		home, _ := os.UserHomeDir()
		eventsPath = filepath.Join(home, "gt", ".events.jsonl")
	}

	events, err := a.readEventsFile(eventsPath, n)
	if err == nil && len(events) > 0 {
		return events, nil
	}

	// Fallback: try gt log command (but it doesn't have JSON output)
	// For now, just return the file-based events or empty
	return events, nil
}

// readEventsFile reads the last N events from a JSONL file.
func (a *Adapter) readEventsFile(path string, n int) ([]model.Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read all lines (for small files this is fine)
	// For large files, we'd want to seek from the end
	var events []model.Event
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var event model.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue // Skip malformed lines
		}
		event.ParsePayload()
		events = append(events, event)
	}

	// Return last N events
	if len(events) > n {
		events = events[len(events)-n:]
	}

	return events, scanner.Err()
}

// StreamEvents returns a channel that emits new events as they occur.
func (a *Adapter) StreamEvents(ctx context.Context) (<-chan model.Event, error) {
	eventsPath := filepath.Join(a.townRoot, ".events.jsonl")
	if a.townRoot == "" {
		home, _ := os.UserHomeDir()
		eventsPath = filepath.Join(home, "gt", ".events.jsonl")
	}

	ch := make(chan model.Event, 100)

	go func() {
		defer close(ch)

		f, err := os.Open(eventsPath)
		if err != nil {
			return
		}
		defer f.Close()

		// Seek to end
		f.Seek(0, 2)

		scanner := bufio.NewScanner(f)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for scanner.Scan() {
					var event model.Event
					if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
						continue
					}
					event.ParsePayload()
					select {
					case ch <- event:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch, nil
}

// GetTownStatus returns the overall town status.
func (a *Adapter) GetTownStatus(ctx context.Context) (*TownStatus, error) {
	cacheKey := "town_status"

	out, err := a.execGT(ctx, "status", "--json")
	if err != nil {
		if cached, ok := a.getCache(cacheKey, 30*time.Second); ok {
			status := cached.(*TownStatus)
			status.Stale = true
			return status, nil
		}
		return nil, err
	}

	var status TownStatus
	if err := parseJSON(out, &status); err != nil {
		return nil, err
	}

	a.setCache(cacheKey, &status)
	return &status, nil
}

// TownStatus represents the overall Gas Town status.
type TownStatus struct {
	Name     string         `json:"name"`
	Path     string         `json:"path"`
	Rigs     []RigStatus    `json:"rigs"`
	Mayor    AgentStatus    `json:"mayor"`
	Deacon   AgentStatus    `json:"deacon"`
	Stale    bool           `json:"-"`
}

// RigStatus represents a rig's status.
type RigStatus struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Prefix   string      `json:"prefix"`
	State    string      `json:"state"` // active, parked, docked
	Witness  AgentStatus `json:"witness"`
	Refinery AgentStatus `json:"refinery"`
}

// AgentStatus represents an agent's running status.
type AgentStatus struct {
	Running   bool   `json:"running"`
	SessionID string `json:"session_id,omitempty"`
}
