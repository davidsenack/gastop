package model

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestConvoyProgress(t *testing.T) {
	c := Convoy{
		TotalCount:  5,
		ClosedCount: 2,
	}
	c.ComputeProgress()

	if c.Progress != 0.4 {
		t.Errorf("expected progress 0.4, got %f", c.Progress)
	}

	if c.ProgressString() != "2/5" {
		t.Errorf("expected '2/5', got '%s'", c.ProgressString())
	}
}

func TestConvoyStatusIcon(t *testing.T) {
	tests := []struct {
		status string
		stuck  bool
		want   string
	}{
		{"closed", false, "✓"},
		{"open", false, "●"},
		{"open", true, "⚠"},
	}

	for _, tt := range tests {
		c := Convoy{Status: tt.status, Stuck: tt.stuck}
		got := c.StatusIcon()
		if got != tt.want {
			t.Errorf("StatusIcon(%s, stuck=%v) = %s, want %s", tt.status, tt.stuck, got, tt.want)
		}
	}
}

func TestBeadStatusIcon(t *testing.T) {
	tests := []struct {
		status string
		stuck  bool
		want   string
	}{
		{"closed", false, "✓"},
		{"in_progress", false, "●"},
		{"blocked", false, "✗"},
		{"deferred", false, "⏸"},
		{"open", false, "○"},
		{"open", true, "⚠"},
	}

	for _, tt := range tests {
		b := Bead{Status: tt.status, Stuck: tt.stuck}
		got := b.StatusIcon()
		if got != tt.want {
			t.Errorf("StatusIcon(%s, stuck=%v) = %s, want %s", tt.status, tt.stuck, got, tt.want)
		}
	}
}

func TestBeadAge(t *testing.T) {
	b := Bead{CreatedAt: time.Now().Add(-2 * time.Hour)}
	b.ComputeAge()

	if b.Age != "2h" {
		t.Errorf("expected age '2h', got '%s'", b.Age)
	}
}

func TestPolecatStateIcon(t *testing.T) {
	tests := []struct {
		state string
		stuck bool
		want  string
	}{
		{"working", false, "●"},
		{"done", false, "✓"},
		{"stuck", false, "⚠"},
		{"idle", false, "○"},
		{"working", true, "⚠"},
	}

	for _, tt := range tests {
		p := Polecat{State: tt.state, Stuck: tt.stuck}
		got := p.StateIcon()
		if got != tt.want {
			t.Errorf("StateIcon(%s, stuck=%v) = %s, want %s", tt.state, tt.stuck, got, tt.want)
		}
	}
}

func TestEventParsing(t *testing.T) {
	data := `{"ts":"2026-01-22T10:00:00Z","source":"gt","type":"spawn","actor":"mayor","payload":{"polecat":"Toast","rig":"gastown"}}`

	var e Event
	if err := json.Unmarshal([]byte(data), &e); err != nil {
		t.Fatalf("failed to parse event: %v", err)
	}

	e.ParsePayload()

	if e.Type != "spawn" {
		t.Errorf("expected type 'spawn', got '%s'", e.Type)
	}
	if e.TargetRig != "gastown" {
		t.Errorf("expected rig 'gastown', got '%s'", e.TargetRig)
	}
	if e.TargetPolecat != "Toast" {
		t.Errorf("expected polecat 'Toast', got '%s'", e.TargetPolecat)
	}
}

func TestParseConvoyFixture(t *testing.T) {
	data, err := os.ReadFile("../../tests/fixtures/convoy_list.json")
	if err != nil {
		t.Skip("fixture file not found")
	}

	var convoys []Convoy
	if err := json.Unmarshal(data, &convoys); err != nil {
		t.Fatalf("failed to parse fixture: %v", err)
	}

	if len(convoys) != 2 {
		t.Errorf("expected 2 convoys, got %d", len(convoys))
	}

	if convoys[0].ID != "hq-abc123" {
		t.Errorf("expected ID 'hq-abc123', got '%s'", convoys[0].ID)
	}
}

func TestParseBeadFixture(t *testing.T) {
	data, err := os.ReadFile("../../tests/fixtures/bead_list.json")
	if err != nil {
		t.Skip("fixture file not found")
	}

	var beads []Bead
	if err := json.Unmarshal(data, &beads); err != nil {
		t.Fatalf("failed to parse fixture: %v", err)
	}

	if len(beads) != 3 {
		t.Errorf("expected 3 beads, got %d", len(beads))
	}

	// Check blocked_by parsing
	if len(beads[2].BlockedBy) != 1 || beads[2].BlockedBy[0] != "gt-002" {
		t.Errorf("expected blocked_by ['gt-002'], got %v", beads[2].BlockedBy)
	}
}
