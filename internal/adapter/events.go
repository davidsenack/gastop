package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

const (
	eventsPollInterval = 100 * time.Millisecond
	eventsChannelSize  = 100
)

// eventsFilePath returns the path to the events JSONL file.
func (a *Adapter) eventsFilePath() string {
	if a.townRoot != "" {
		return filepath.Join(a.townRoot, ".events.jsonl")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "gt", ".events.jsonl")
}

// parseEvent parses a JSON line into an Event struct.
func parseEvent(data []byte) (model.Event, bool) {
	var event model.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return event, false
	}
	event.ParsePayload()
	return event, true
}

// TailEvents returns the last N events from the activity log.
func (a *Adapter) TailEvents(ctx context.Context, n int) ([]model.Event, error) {
	events, err := a.readEventsFile(a.eventsFilePath(), n)
	if err == nil && len(events) > 0 {
		return events, nil
	}
	return events, nil
}

// readEventsFile reads the last N events from a JSONL file.
func (a *Adapter) readEventsFile(path string, n int) ([]model.Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var events []model.Event
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if event, ok := parseEvent(scanner.Bytes()); ok {
			events = append(events, event)
		}
	}

	if len(events) > n {
		events = events[len(events)-n:]
	}
	return events, scanner.Err()
}

// StreamEvents returns a channel that emits new events as they occur.
func (a *Adapter) StreamEvents(ctx context.Context) (<-chan model.Event, error) {
	ch := make(chan model.Event, eventsChannelSize)
	go a.streamEventsLoop(ctx, ch)
	return ch, nil
}

// streamEventsLoop handles the event streaming goroutine.
func (a *Adapter) streamEventsLoop(ctx context.Context, ch chan<- model.Event) {
	defer close(ch)

	f, err := os.Open(a.eventsFilePath())
	if err != nil {
		return
	}
	defer f.Close()

	// Seek to end to only receive new events
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	ticker := time.NewTicker(eventsPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.processNewEvents(ctx, scanner, ch)
		}
	}
}

// processNewEvents reads and sends any new events from the scanner.
func (a *Adapter) processNewEvents(ctx context.Context, scanner *bufio.Scanner, ch chan<- model.Event) {
	for scanner.Scan() {
		event, ok := parseEvent(scanner.Bytes())
		if !ok {
			continue
		}
		select {
		case ch <- event:
		case <-ctx.Done():
			return
		}
	}
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
