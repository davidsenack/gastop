package model

import (
	"fmt"
	"time"
)

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
	Branch       string    `json:"branch,omitempty"`
	ClonePath    string    `json:"clone_path,omitempty"`
	Windows      int       `json:"windows,omitempty"`

	// Hooked work info (populated separately)
	HookedBead  string `json:"-"`
	HookedTitle string `json:"-"`

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

// ActivityAgo returns a human-readable string of time since last activity.
func (p *Polecat) ActivityAgo() string {
	if p.LastActivity.IsZero() {
		return ""
	}
	d := time.Since(p.LastActivity)
	if d < time.Minute {
		return "now"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		return fmt.Sprintf("%dm", mins)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		return fmt.Sprintf("%dh", hours)
	}
	days := int(d.Hours() / 24)
	return fmt.Sprintf("%dd", days)
}

// WorkDescription returns a short description of current work.
func (p *Polecat) WorkDescription() string {
	if p.HookedBead != "" {
		if p.HookedTitle != "" {
			// Truncate title if too long
			title := p.HookedTitle
			if len(title) > 30 {
				title = title[:27] + "..."
			}
			return p.HookedBead + ": " + title
		}
		return p.HookedBead
	}
	if p.AssignedBead != "" {
		return p.AssignedBead
	}
	if p.Branch != "" {
		// Extract meaningful part of branch name
		branch := p.Branch
		if len(branch) > 25 {
			branch = "..." + branch[len(branch)-22:]
		}
		return branch
	}
	return ""
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
