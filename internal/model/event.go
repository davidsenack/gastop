package model

import (
	"encoding/json"
	"time"
)

// eventIcons maps event types to their display icons.
var eventIcons = map[string]string{
	"session_start":  "▶",
	"spawn":          "+",
	"sling":          "🎯",
	"handoff":        "🤝",
	"done":           "✓",
	"crash":          "💥",
	"kill":           "✗",
	"nudge":          "⚡",
	"polecat_nudged": "⚡",
	"patrol_started": "🦉",
	"merge_started":  "⚙",
	"merged":         "✓",
	"merge_failed":   "✗",
	"create":         "+",
	"bonded":         "+",
	"update":         "~",
	"delete":         "⊘",
	"in_progress":    "→",
	"completed":      "✓",
	"failed":         "✗",
}

// eventSummaries maps event types to their summary templates.
var eventSummaries = map[string]string{
	"session_start": "session started",
	"handoff":       "handed off",
	"done":          "completed",
	"crash":         "crashed",
	"kill":          "killed",
}

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
	if icon, ok := eventIcons[e.Type]; ok {
		return icon
	}
	return "•"
}

// Summary returns a short description of the event.
func (e *Event) Summary() string {
	// Handle special cases with dynamic content
	switch e.Type {
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
	}

	// Use static mapping for simple cases
	if summary, ok := eventSummaries[e.Type]; ok {
		return summary
	}
	return e.Type
}

// TimeString returns a short time string.
func (e *Event) TimeString() string {
	return e.Timestamp.Format("15:04:05")
}
