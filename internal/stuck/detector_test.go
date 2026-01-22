package stuck

import (
	"testing"
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

func TestDetectorCheckBeads(t *testing.T) {
	d := NewDetector(30)

	beads := []model.Bead{
		{
			ID:        "test-1",
			Status:    "in_progress",
			UpdatedAt: time.Now().Add(-45 * time.Minute), // Stuck - no updates for 45min
		},
		{
			ID:        "test-2",
			Status:    "in_progress",
			UpdatedAt: time.Now().Add(-10 * time.Minute), // Not stuck - recent update
		},
		{
			ID:        "test-3",
			Status:    "closed",
			UpdatedAt: time.Now().Add(-2 * time.Hour), // Not stuck - closed
		},
	}

	d.CheckBeads(beads)

	if !beads[0].Stuck {
		t.Error("expected bead test-1 to be stuck")
	}
	if beads[0].StuckReason == "" {
		t.Error("expected stuck reason for test-1")
	}

	if beads[1].Stuck {
		t.Error("expected bead test-2 to not be stuck")
	}

	if beads[2].Stuck {
		t.Error("expected bead test-3 to not be stuck")
	}
}

func TestDetectorCheckPolecats(t *testing.T) {
	d := NewDetector(30)

	polecats := []model.Polecat{
		{
			Name:         "Toast",
			State:        "working",
			AssignedBead: "", // Stuck - working but no assigned bead
			Running:      true,
		},
		{
			Name:         "Furiosa",
			State:        "working",
			AssignedBead: "gt-001",
			Running:      false, // Stuck - session not running
		},
		{
			Name:         "Nux",
			State:        "idle",
			AssignedBead: "",
			Running:      true, // Not stuck - idle is fine
		},
		{
			Name:         "Max",
			State:        "stuck", // Already marked stuck by GT
			AssignedBead: "gt-002",
			Running:      true,
		},
	}

	d.CheckPolecats(polecats)

	if !polecats[0].Stuck {
		t.Error("expected Toast to be stuck (working without bead)")
	}
	if polecats[0].StuckReason != "Working but no assigned bead" {
		t.Errorf("unexpected stuck reason: %s", polecats[0].StuckReason)
	}

	if !polecats[1].Stuck {
		t.Error("expected Furiosa to be stuck (session not running)")
	}

	if polecats[2].Stuck {
		t.Error("expected Nux to not be stuck")
	}

	if !polecats[3].Stuck {
		t.Error("expected Max to be stuck (already marked)")
	}
}

func TestDetectorSummarize(t *testing.T) {
	d := NewDetector(30)

	beads := []model.Bead{
		{Stuck: true},
		{Stuck: true},
		{Stuck: false},
	}
	polecats := []model.Polecat{
		{Stuck: true},
	}
	convoys := []model.Convoy{
		{Stuck: false},
		{Stuck: true},
	}

	summary := d.Summarize(beads, polecats, convoys)

	if summary.StuckBeads != 2 {
		t.Errorf("expected 2 stuck beads, got %d", summary.StuckBeads)
	}
	if summary.StuckPolecats != 1 {
		t.Errorf("expected 1 stuck polecat, got %d", summary.StuckPolecats)
	}
	if summary.StuckConvoys != 1 {
		t.Errorf("expected 1 stuck convoy, got %d", summary.StuckConvoys)
	}
}
