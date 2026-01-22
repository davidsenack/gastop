package model

import (
	"fmt"
	"time"
)

// Convoy represents a batch of tracked work across rigs.
// Supports both gt convoy list format and bd list -t convoy format.
type Convoy struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"` // open, closed
	TrackedIDs  []string  `json:"tracked_ids,omitempty"`
	TotalCount  int       `json:"total_count"`
	ClosedCount int       `json:"closed_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	Owner       string    `json:"owner,omitempty"`

	// Fields from bd list -t convoy format
	Description     string `json:"description,omitempty"`
	Priority        int    `json:"priority,omitempty"`
	IssueType       string `json:"issue_type,omitempty"`
	DependencyCount int    `json:"dependency_count,omitempty"`

	// Computed fields
	Progress    float64 `json:"-"` // 0.0 - 1.0
	Stuck       bool    `json:"-"`
	StuckReason string  `json:"-"`
}

// ComputeProgress calculates the completion percentage.
func (c *Convoy) ComputeProgress() {
	// Use DependencyCount as TotalCount if not set (bd format)
	if c.TotalCount == 0 && c.DependencyCount > 0 {
		c.TotalCount = c.DependencyCount
	}
	if c.TotalCount > 0 {
		c.Progress = float64(c.ClosedCount) / float64(c.TotalCount)
	}
}

// ProgressString returns a human-readable progress string.
func (c *Convoy) ProgressString() string {
	return fmt.Sprintf("%d/%d", c.ClosedCount, c.TotalCount)
}

// StatusIcon returns a status indicator character.
func (c *Convoy) StatusIcon() string {
	switch c.Status {
	case "closed":
		return "✓"
	case "open":
		if c.Stuck {
			return "⚠"
		}
		return "●"
	default:
		return "○"
	}
}
