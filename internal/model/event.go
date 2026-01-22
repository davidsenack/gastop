package model

import (
	"encoding/json"
	"time"
)

// Event represents an activity event from Gas Town.
type Event struct {
	Timestamp  time.Time       `json:"ts"`
	Source     string          `json:"source"` // gt, bd
	Type       string          `json:"type"`   // session_start, spawn, sling, handoff, etc.
	Actor      string          `json:"actor"`
	Payload    json.RawMessage `json:"payload,omitempty"`
	Visibility string          `json:"visibility,omitempty"`

	// Parsed payload fields (type-dependent)
	TargetRig     string `json:"-"`
	TargetPolecat string `json:"-"`
	TargetBead    string `json:"-"`
	Message       string `json:"-"`
}

// ParsePayload extracts common fields from the payload based on event type.
func (e *Event) ParsePayload() {
	if len(e.Payload) == 0 {
		return
	}

	var p map[string]interface{}
	if err := json.Unmarshal(e.Payload, &p); err != nil {
		return
	}

	// Extract common fields
	if v, ok := p["rig"].(string); ok {
		e.TargetRig = v
	}
	if v, ok := p["polecat"].(string); ok {
		e.TargetPolecat = v
	}
	if v, ok := p["bead"].(string); ok {
		e.TargetBead = v
	}
	if v, ok := p["target"].(string); ok {
		// target can contain rig/polecat path
		e.Message = v
	}
}

// Icon returns an icon for the event type.
func (e *Event) Icon() string {
	switch e.Type {
	case "session_start":
		return "▶"
	case "spawn":
		return "+"
	case "sling":
		return "🎯"
	case "handoff":
		return "🤝"
	case "done":
		return "✓"
	case "crash":
		return "💥"
	case "kill":
		return "✗"
	case "nudge", "polecat_nudged":
		return "⚡"
	case "patrol_started":
		return "🦉"
	case "merge_started":
		return "⚙"
	case "merged":
		return "✓"
	case "merge_failed":
		return "✗"
	case "create", "bonded":
		return "+"
	case "update":
		return "~"
	case "delete":
		return "⊘"
	case "in_progress":
		return "→"
	case "completed":
		return "✓"
	case "failed":
		return "✗"
	default:
		return "•"
	}
}

// Summary returns a short description of the event.
func (e *Event) Summary() string {
	switch e.Type {
	case "session_start":
		return "session started"
	case "spawn":
		if e.TargetPolecat != "" {
			return "spawned " + e.TargetPolecat
		}
		return "spawned"
	case "sling":
		if e.TargetBead != "" && e.Message != "" {
			return e.TargetBead + " → " + e.Message
		}
		return "slung work"
	case "handoff":
		return "handed off"
	case "done":
		return "completed"
	case "crash":
		return "crashed"
	case "kill":
		return "killed"
	default:
		return e.Type
	}
}

// TimeString returns a short time string.
func (e *Event) TimeString() string {
	return e.Timestamp.Format("15:04:05")
}
