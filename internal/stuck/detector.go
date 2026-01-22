package stuck

import (
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

// Detector checks for stuck work.
type Detector struct {
	// StuckThreshold is how long a bead can be in_progress without updates
	// before being considered stuck.
	StuckThreshold time.Duration

	// HeartbeatThreshold is how long since a polecat's last activity
	// before being considered stuck.
	HeartbeatThreshold time.Duration
}

// NewDetector creates a detector with default thresholds.
func NewDetector(stuckMinutes int) *Detector {
	if stuckMinutes <= 0 {
		stuckMinutes = 30
	}
	return &Detector{
		StuckThreshold:     time.Duration(stuckMinutes) * time.Minute,
		HeartbeatThreshold: time.Duration(stuckMinutes) * time.Minute,
	}
}

// CheckBeads marks beads as stuck if they meet stuck criteria.
func (d *Detector) CheckBeads(beads []model.Bead) {
	for i := range beads {
		d.checkBead(&beads[i])
	}
}

// checkBead checks a single bead for stuck conditions.
func (d *Detector) checkBead(b *model.Bead) {
	// Already closed or deferred - not stuck
	if b.Status == "closed" || b.Status == "deferred" {
		return
	}

	// In progress but no updates for too long
	if b.Status == "in_progress" {
		if time.Since(b.UpdatedAt) > d.StuckThreshold {
			b.Stuck = true
			b.StuckReason = "No updates for " + humanizeDuration(time.Since(b.UpdatedAt))
			return
		}
	}

	// Blocked for a long time
	if b.Status == "blocked" {
		if time.Since(b.UpdatedAt) > d.StuckThreshold*2 {
			b.Stuck = true
			b.StuckReason = "Blocked for " + humanizeDuration(time.Since(b.UpdatedAt))
			return
		}
	}
}

// CheckPolecats marks polecats as stuck if they meet stuck criteria.
func (d *Detector) CheckPolecats(polecats []model.Polecat) {
	for i := range polecats {
		d.checkPolecat(&polecats[i])
	}
}

// checkPolecat checks a single polecat for stuck conditions.
func (d *Detector) checkPolecat(p *model.Polecat) {
	// Already marked stuck by Gas Town
	if p.State == "stuck" {
		p.Stuck = true
		if p.StuckReason == "" {
			p.StuckReason = "Marked stuck by Gas Town"
		}
		return
	}

	// Working but no assigned bead
	if p.State == "working" && p.AssignedBead == "" {
		p.Stuck = true
		p.StuckReason = "Working but no assigned bead"
		return
	}

	// Working but no activity for too long
	if p.State == "working" && !p.LastActivity.IsZero() {
		if time.Since(p.LastActivity) > d.HeartbeatThreshold {
			p.Stuck = true
			p.StuckReason = "No activity for " + humanizeDuration(time.Since(p.LastActivity))
			return
		}
	}

	// Session not running but state says working
	if p.State == "working" && !p.Running {
		p.Stuck = true
		p.StuckReason = "Session not running"
		return
	}
}

// CheckConvoys marks convoys as stuck if they have stuck beads.
func (d *Detector) CheckConvoys(convoys []model.Convoy, getBeads func(ids []string) []model.Bead) {
	for i := range convoys {
		d.checkConvoy(&convoys[i], getBeads)
	}
}

// checkConvoy checks a convoy for stuck conditions.
func (d *Detector) checkConvoy(c *model.Convoy, getBeads func(ids []string) []model.Bead) {
	if c.Status == "closed" {
		return
	}

	// Check if any tracked beads are stuck
	if getBeads != nil && len(c.TrackedIDs) > 0 {
		beads := getBeads(c.TrackedIDs)
		for _, b := range beads {
			if b.Stuck {
				c.Stuck = true
				c.StuckReason = "Has stuck beads"
				return
			}
		}
	}

	// Determine the reference time - use UpdatedAt if set, otherwise CreatedAt
	refTime := c.UpdatedAt
	if refTime.IsZero() {
		refTime = c.CreatedAt
	}

	// Skip stuck check if we still don't have a valid timestamp
	if refTime.IsZero() {
		return
	}

	// No progress for a long time
	if time.Since(refTime) > d.StuckThreshold*2 {
		c.Stuck = true
		c.StuckReason = "No progress for " + humanizeDuration(time.Since(refTime))
	}
}

// StuckSummary returns counts of stuck items.
type StuckSummary struct {
	StuckBeads    int
	StuckPolecats int
	StuckConvoys  int
}

// Summarize returns a summary of stuck items.
func (d *Detector) Summarize(beads []model.Bead, polecats []model.Polecat, convoys []model.Convoy) StuckSummary {
	var s StuckSummary
	for _, b := range beads {
		if b.Stuck {
			s.StuckBeads++
		}
	}
	for _, p := range polecats {
		if p.Stuck {
			s.StuckPolecats++
		}
	}
	for _, c := range convoys {
		if c.Stuck {
			s.StuckConvoys++
		}
	}
	return s
}

func humanizeDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}
	if d < 24*time.Hour {
		return d.Round(time.Hour).String()
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return string(rune(days+'0')) + " days"
}
