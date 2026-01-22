package model

import "time"

// Polecat represents a worker agent in Gas Town.
type Polecat struct {
	Name         string    `json:"name"`
	Rig          string    `json:"rig"`
	State        string    `json:"state"` // working, done, stuck, idle
	AssignedBead string    `json:"assigned_bead,omitempty"`
	SessionID    string    `json:"session_id,omitempty"`
	Running      bool      `json:"session_running"`
	Attached     bool      `json:"attached"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`

	// Computed fields
	Stuck       bool   `json:"-"`
	StuckReason string `json:"-"`
}

// FullName returns the rig/name format.
func (p *Polecat) FullName() string {
	return p.Rig + "/" + p.Name
}

// StateIcon returns a state indicator character.
func (p *Polecat) StateIcon() string {
	if p.Stuck {
		return "⚠"
	}
	switch p.State {
	case "working":
		return "●"
	case "done":
		return "✓"
	case "stuck":
		return "⚠"
	case "idle":
		return "○"
	default:
		return "?"
	}
}

// SessionStatus returns a human-readable session status.
func (p *Polecat) SessionStatus() string {
	if !p.Running {
		return "stopped"
	}
	if p.Attached {
		return "attached"
	}
	return "running"
}

// TimeSinceActivity returns duration since last activity.
func (p *Polecat) TimeSinceActivity() time.Duration {
	return time.Since(p.LastActivity)
}

// Agent represents a broader category of agents (witness, refinery, crew).
type Agent struct {
	Name      string `json:"name"`
	Role      string `json:"role"` // witness, refinery, crew, polecat
	Rig       string `json:"rig"`
	State     string `json:"state"`
	Running   bool   `json:"running"`
	SessionID string `json:"session_id,omitempty"`
}

// FullName returns the rig/name format.
func (a *Agent) FullName() string {
	if a.Rig == "" {
		return a.Name
	}
	return a.Rig + "/" + a.Name
}
