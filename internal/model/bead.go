package model

import (
	"fmt"
	"time"
)

// Bead represents an issue/task in the beads system.
type Bead struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"` // open, in_progress, blocked, deferred, closed
	Priority    int        `json:"priority"` // 0-4, 0=highest
	IssueType   string     `json:"issue_type"` // bug, feature, task, epic, molecule, agent
	Owner       string     `json:"owner,omitempty"`
	Assignee    string     `json:"assignee,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	CloseReason string     `json:"close_reason,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
	Parent      string     `json:"parent,omitempty"`
	Ephemeral   bool       `json:"ephemeral,omitempty"`

	// Dependency info (from bd show)
	Blocks    []string `json:"blocks,omitempty"`
	BlockedBy []string `json:"blocked_by,omitempty"`

	// Computed fields
	Stuck       bool   `json:"-"`
	StuckReason string `json:"-"`
	Age         string `json:"-"` // Human-readable age
}

// StatusIcon returns a status indicator character.
func (b *Bead) StatusIcon() string {
	if b.Stuck {
		return "⚠"
	}
	switch b.Status {
	case "closed":
		return "✓"
	case "in_progress":
		return "●"
	case "blocked":
		return "✗"
	case "deferred":
		return "⏸"
	case "open":
		return "○"
	default:
		return "?"
	}
}

// PriorityString returns P0-P4 format.
func (b *Bead) PriorityString() string {
	return fmt.Sprintf("P%d", b.Priority)
}

// ComputeAge sets the human-readable age field.
func (b *Bead) ComputeAge() {
	b.Age = humanizeDuration(time.Since(b.CreatedAt))
}

// TimeSinceUpdate returns duration since last update.
func (b *Bead) TimeSinceUpdate() time.Duration {
	return time.Since(b.UpdatedAt)
}

// IsBlocked returns true if the bead has blockers.
func (b *Bead) IsBlocked() bool {
	return len(b.BlockedBy) > 0 || b.Status == "blocked"
}

// DependencyCount returns the number of dependencies.
func (b *Bead) DependencyCount() int {
	return len(b.BlockedBy)
}

// humanizeDuration converts a duration to a short human-readable string.
func humanizeDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	return fmt.Sprintf("%dd", days)
}
