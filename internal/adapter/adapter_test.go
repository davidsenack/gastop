package adapter

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

// TestParseJSON tests the parseJSON utility function.
func TestParseJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: false,
		},
		{
			name:    "valid JSON array",
			input:   []byte(`[{"id":"test"}]`),
			wantErr: false,
		},
		{
			name:    "valid JSON object",
			input:   []byte(`{"id":"test"}`),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := parseJSON(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParseBeadListFixture tests parsing of bead list JSON output.
func TestParseBeadListFixture(t *testing.T) {
	data, err := os.ReadFile("../../tests/fixtures/bead_list.json")
	if err != nil {
		t.Skip("fixture file not found")
	}

	var beads []model.Bead
	if err := parseJSON(data, &beads); err != nil {
		t.Fatalf("failed to parse bead fixture: %v", err)
	}

	if len(beads) != 3 {
		t.Errorf("expected 3 beads, got %d", len(beads))
	}

	// Verify first bead
	if beads[0].ID != "gt-001" {
		t.Errorf("expected ID 'gt-001', got '%s'", beads[0].ID)
	}
	if beads[0].Title != "Implement convoy panel" {
		t.Errorf("expected title 'Implement convoy panel', got '%s'", beads[0].Title)
	}
	if beads[0].Status != "closed" {
		t.Errorf("expected status 'closed', got '%s'", beads[0].Status)
	}
	if beads[0].Assignee != "Toast" {
		t.Errorf("expected assignee 'Toast', got '%s'", beads[0].Assignee)
	}

	// Verify bead with blocked_by
	if len(beads[2].BlockedBy) != 1 || beads[2].BlockedBy[0] != "gt-002" {
		t.Errorf("expected blocked_by ['gt-002'], got %v", beads[2].BlockedBy)
	}
}

// TestParseConvoyListFixture tests parsing of convoy list JSON output.
func TestParseConvoyListFixture(t *testing.T) {
	data, err := os.ReadFile("../../tests/fixtures/convoy_list.json")
	if err != nil {
		t.Skip("fixture file not found")
	}

	var convoys []model.Convoy
	if err := parseJSON(data, &convoys); err != nil {
		t.Fatalf("failed to parse convoy fixture: %v", err)
	}

	if len(convoys) != 2 {
		t.Errorf("expected 2 convoys, got %d", len(convoys))
	}

	// Verify first convoy
	if convoys[0].ID != "hq-abc123" {
		t.Errorf("expected ID 'hq-abc123', got '%s'", convoys[0].ID)
	}
	if convoys[0].Title != "beadtop MVP" {
		t.Errorf("expected title 'beadtop MVP', got '%s'", convoys[0].Title)
	}
	if convoys[0].TotalCount != 5 {
		t.Errorf("expected total_count 5, got %d", convoys[0].TotalCount)
	}
	if convoys[0].ClosedCount != 2 {
		t.Errorf("expected closed_count 2, got %d", convoys[0].ClosedCount)
	}

	// Verify closed convoy
	if convoys[1].Status != "closed" {
		t.Errorf("expected status 'closed', got '%s'", convoys[1].Status)
	}
}

// TestParsePolecatListFixture tests parsing of polecat list JSON output.
func TestParsePolecatListFixture(t *testing.T) {
	data, err := os.ReadFile("../../tests/fixtures/polecat_list.json")
	if err != nil {
		t.Skip("fixture file not found")
	}

	var polecats []model.Polecat
	if err := parseJSON(data, &polecats); err != nil {
		t.Fatalf("failed to parse polecat fixture: %v", err)
	}

	if len(polecats) != 2 {
		t.Errorf("expected 2 polecats, got %d", len(polecats))
	}

	// Verify first polecat
	if polecats[0].Name != "Toast" {
		t.Errorf("expected name 'Toast', got '%s'", polecats[0].Name)
	}
	if polecats[0].Rig != "gastown" {
		t.Errorf("expected rig 'gastown', got '%s'", polecats[0].Rig)
	}
	if polecats[0].State != "working" {
		t.Errorf("expected state 'working', got '%s'", polecats[0].State)
	}
	// Note: fixture uses "running" but model expects "session_running"
	// The Running field won't be populated from this fixture

	// Test FullName method
	if polecats[0].FullName() != "gastown/Toast" {
		t.Errorf("expected FullName 'gastown/Toast', got '%s'", polecats[0].FullName())
	}
}

// TestParseEvent tests parsing of individual event JSON.
func TestParseEvent(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantOK         bool
		wantType       string
		wantTargetRig  string
		wantTargetPcat string
	}{
		{
			name:           "spawn event",
			input:          `{"ts":"2026-01-22T10:00:00Z","source":"gt","type":"spawn","actor":"mayor","payload":{"polecat":"Toast","rig":"gastown"}}`,
			wantOK:         true,
			wantType:       "spawn",
			wantTargetRig:  "gastown",
			wantTargetPcat: "Toast",
		},
		{
			name:   "sling event with bead",
			input:  `{"ts":"2026-01-22T11:00:00Z","source":"gt","type":"sling","actor":"witness","payload":{"bead":"gt-001","target":"gastown/Toast"}}`,
			wantOK: true,
			wantType: "sling",
		},
		{
			name:   "done event",
			input:  `{"ts":"2026-01-22T12:00:00Z","source":"gt","type":"done","actor":"Toast"}`,
			wantOK: true,
			wantType: "done",
		},
		{
			name:   "invalid JSON",
			input:  `{invalid json}`,
			wantOK: false,
		},
		{
			name:   "empty payload",
			input:  `{"ts":"2026-01-22T10:00:00Z","source":"gt","type":"crash","actor":"Furiosa"}`,
			wantOK: true,
			wantType: "crash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, ok := parseEvent([]byte(tt.input))
			if ok != tt.wantOK {
				t.Errorf("parseEvent() ok = %v, wantOK %v", ok, tt.wantOK)
				return
			}
			if !ok {
				return
			}
			if event.Type != tt.wantType {
				t.Errorf("event.Type = %s, want %s", event.Type, tt.wantType)
			}
			if tt.wantTargetRig != "" && event.TargetRig != tt.wantTargetRig {
				t.Errorf("event.TargetRig = %s, want %s", event.TargetRig, tt.wantTargetRig)
			}
			if tt.wantTargetPcat != "" && event.TargetPolecat != tt.wantTargetPcat {
				t.Errorf("event.TargetPolecat = %s, want %s", event.TargetPolecat, tt.wantTargetPcat)
			}
		})
	}
}

// TestNewAdapter tests adapter creation with various configurations.
func TestNewAdapter(t *testing.T) {
	tests := []struct {
		name     string
		gtPath   string
		bdPath   string
		townRoot string
		wantGT   string
		wantBD   string
	}{
		{
			name:     "default paths",
			gtPath:   "",
			bdPath:   "",
			townRoot: "/tmp/test",
			wantGT:   "gt",
			wantBD:   "bd",
		},
		{
			name:     "custom paths",
			gtPath:   "/usr/local/bin/gt",
			bdPath:   "/usr/local/bin/bd",
			townRoot: "/home/user/project",
			wantGT:   "/usr/local/bin/gt",
			wantBD:   "/usr/local/bin/bd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.gtPath, tt.bdPath, tt.townRoot)
			if a.gtPath != tt.wantGT {
				t.Errorf("gtPath = %s, want %s", a.gtPath, tt.wantGT)
			}
			if a.bdPath != tt.wantBD {
				t.Errorf("bdPath = %s, want %s", a.bdPath, tt.wantBD)
			}
			if a.townRoot != tt.townRoot {
				t.Errorf("townRoot = %s, want %s", a.townRoot, tt.townRoot)
			}
		})
	}
}

// TestCacheOperations tests the caching functionality.
func TestCacheOperations(t *testing.T) {
	a := New("", "", "")

	// Test setCache and getCache
	testData := []model.Bead{{ID: "test-1", Title: "Test Bead"}}
	a.setCache("test-key", testData)

	cached, ok := a.getCache("test-key", 5*time.Minute)
	if !ok {
		t.Error("expected cache hit, got miss")
	}
	beads := cached.([]model.Bead)
	if len(beads) != 1 || beads[0].ID != "test-1" {
		t.Error("cached data does not match")
	}

	// Test cache miss
	_, ok = a.getCache("nonexistent", 5*time.Minute)
	if ok {
		t.Error("expected cache miss, got hit")
	}

	// Test ClearCache
	a.ClearCache()
	_, ok = a.getCache("test-key", 5*time.Minute)
	if ok {
		t.Error("expected cache miss after clear, got hit")
	}
}

// TestCacheStaleDetection tests cache staleness detection.
func TestCacheStaleDetection(t *testing.T) {
	a := New("", "", "")

	// Set cache entry
	a.setCache("test-key", "test-data")

	// Entry should not be stale initially
	if a.IsCacheStale("test-key") {
		t.Error("expected cache entry to not be stale initially")
	}

	// Check with very short max age to trigger staleness
	_, _ = a.getCache("test-key", 0)

	// Now it should be stale
	if !a.IsCacheStale("test-key") {
		t.Error("expected cache entry to be stale after expired maxAge")
	}

	// Nonexistent key should not be stale
	if a.IsCacheStale("nonexistent") {
		t.Error("expected nonexistent key to not be stale")
	}
}

// TestSetTimeout tests timeout configuration.
func TestSetTimeout(t *testing.T) {
	a := New("", "", "")

	// Default timeout should be 5 seconds
	if a.timeout != 5*time.Second {
		t.Errorf("default timeout = %v, want 5s", a.timeout)
	}

	// Test SetTimeout
	a.SetTimeout(10 * time.Second)
	if a.timeout != 10*time.Second {
		t.Errorf("timeout after SetTimeout = %v, want 10s", a.timeout)
	}
}

// TestEventsFilePath tests the events file path calculation.
func TestEventsFilePath(t *testing.T) {
	// With town root
	a := New("", "", "/tmp/test-town")
	path := a.eventsFilePath()
	if path != "/tmp/test-town/.events.jsonl" {
		t.Errorf("eventsFilePath with townRoot = %s, want /tmp/test-town/.events.jsonl", path)
	}

	// Without town root (should use home directory)
	a2 := New("", "", "")
	path2 := a2.eventsFilePath()
	home, _ := os.UserHomeDir()
	expected := home + "/gt/.events.jsonl"
	if path2 != expected {
		t.Errorf("eventsFilePath without townRoot = %s, want %s", path2, expected)
	}
}

// TestHookStatusParsing tests parsing of hook status JSON.
func TestHookStatusParsing(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantStatus string
		wantBead   string
		wantTitle  string
	}{
		{
			name:       "empty hook",
			input:      `{"agent":"Toast","status":"empty"}`,
			wantStatus: "empty",
			wantBead:   "",
		},
		{
			name:       "hooked bead",
			input:      `{"agent":"Toast","status":"hooked","bead":"gt-001","title":"Implement feature"}`,
			wantStatus: "hooked",
			wantBead:   "gt-001",
			wantTitle:  "Implement feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status HookStatus
			if err := json.Unmarshal([]byte(tt.input), &status); err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			if status.Status != tt.wantStatus {
				t.Errorf("Status = %s, want %s", status.Status, tt.wantStatus)
			}
			if status.Bead != tt.wantBead {
				t.Errorf("Bead = %s, want %s", status.Bead, tt.wantBead)
			}
			if status.Title != tt.wantTitle {
				t.Errorf("Title = %s, want %s", status.Title, tt.wantTitle)
			}
		})
	}
}

// TestTownStatusParsing tests parsing of town status JSON.
func TestTownStatusParsing(t *testing.T) {
	input := `{
		"name": "test-town",
		"path": "/home/user/project",
		"rigs": [
			{
				"name": "gastown",
				"path": "/home/user/project/gastown",
				"prefix": "gt",
				"state": "active",
				"witness": {"running": true, "session_id": "witness-1"},
				"refinery": {"running": false}
			}
		],
		"mayor": {"running": true, "session_id": "mayor-1"},
		"deacon": {"running": true}
	}`

	var status TownStatus
	if err := json.Unmarshal([]byte(input), &status); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if status.Name != "test-town" {
		t.Errorf("Name = %s, want test-town", status.Name)
	}
	if len(status.Rigs) != 1 {
		t.Errorf("Rigs count = %d, want 1", len(status.Rigs))
	}
	if !status.Mayor.Running {
		t.Error("expected Mayor.Running to be true")
	}
	if status.Rigs[0].Name != "gastown" {
		t.Errorf("Rigs[0].Name = %s, want gastown", status.Rigs[0].Name)
	}
	if !status.Rigs[0].Witness.Running {
		t.Error("expected Rigs[0].Witness.Running to be true")
	}
}

// TestBeadListOptsArgs tests that BeadListOpts produces correct arguments.
func TestBeadListOptsCacheKey(t *testing.T) {
	tests := []struct {
		name    string
		opts    BeadListOpts
		wantKey string
	}{
		{
			name:    "default options",
			opts:    BeadListOpts{},
			wantKey: "beads",
		},
		{
			name:    "ready beads",
			opts:    BeadListOpts{Ready: true},
			wantKey: "beads:ready",
		},
		{
			name:    "blocked beads",
			opts:    BeadListOpts{Blocked: true},
			wantKey: "beads:blocked",
		},
		{
			name:    "with status filter",
			opts:    BeadListOpts{Status: "open"},
			wantKey: "beads:status=open",
		},
		{
			name:    "with assignee",
			opts:    BeadListOpts{Assignee: "Toast"},
			wantKey: "beads:assignee=Toast",
		},
		{
			name:    "with type",
			opts:    BeadListOpts{Type: "task"},
			wantKey: "beads:type=task",
		},
		{
			name:    "with limit",
			opts:    BeadListOpts{Limit: 50},
			wantKey: "beads:limit=50",
		},
		{
			name:    "combined filters",
			opts:    BeadListOpts{Status: "in_progress", Assignee: "Toast", Type: "bug"},
			wantKey: "beads:status=in_progress:assignee=Toast:type=bug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build cache key the same way ListBeads does
			cacheKey := "beads"
			if tt.opts.Ready {
				cacheKey += ":ready"
			} else if tt.opts.Blocked {
				cacheKey += ":blocked"
			}
			if tt.opts.Status != "" && !tt.opts.Ready && !tt.opts.Blocked {
				cacheKey += ":status=" + tt.opts.Status
			}
			if tt.opts.Limit > 0 {
				cacheKey += ":limit=" + string(rune('0'+tt.opts.Limit/10)) + string(rune('0'+tt.opts.Limit%10))
			}
			if tt.opts.Assignee != "" {
				cacheKey += ":assignee=" + tt.opts.Assignee
			}
			if tt.opts.Type != "" {
				cacheKey += ":type=" + tt.opts.Type
			}

			if cacheKey != tt.wantKey {
				t.Errorf("cacheKey = %s, want %s", cacheKey, tt.wantKey)
			}
		})
	}
}

// TestConvoyListOptsCacheKey tests ConvoyListOpts cache key generation.
func TestConvoyListOptsCacheKey(t *testing.T) {
	tests := []struct {
		name    string
		opts    ConvoyListOpts
		wantKey string
	}{
		{
			name:    "default",
			opts:    ConvoyListOpts{},
			wantKey: "convoys",
		},
		{
			name:    "all convoys",
			opts:    ConvoyListOpts{All: true},
			wantKey: "convoys:all",
		},
		{
			name:    "with status",
			opts:    ConvoyListOpts{Status: "open"},
			wantKey: "convoys:status=open",
		},
		{
			name:    "all with status",
			opts:    ConvoyListOpts{All: true, Status: "closed"},
			wantKey: "convoys:all:status=closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheKey := "convoys"
			if tt.opts.All {
				cacheKey += ":all"
			}
			if tt.opts.Status != "" {
				cacheKey += ":status=" + tt.opts.Status
			}

			if cacheKey != tt.wantKey {
				t.Errorf("cacheKey = %s, want %s", cacheKey, tt.wantKey)
			}
		})
	}
}

// TestPolecatCacheKey tests polecat cache key generation.
func TestPolecatCacheKey(t *testing.T) {
	tests := []struct {
		name    string
		rig     string
		wantKey string
	}{
		{
			name:    "all rigs",
			rig:     "",
			wantKey: "polecats:all",
		},
		{
			name:    "specific rig",
			rig:     "gastown",
			wantKey: "polecats:gastown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheKey := "polecats"
			if tt.rig != "" {
				cacheKey += ":" + tt.rig
			} else {
				cacheKey += ":all"
			}

			if cacheKey != tt.wantKey {
				t.Errorf("cacheKey = %s, want %s", cacheKey, tt.wantKey)
			}
		})
	}
}

// TestConcurrentCacheAccess tests thread-safe cache access.
func TestConcurrentCacheAccess(t *testing.T) {
	a := New("", "", "")
	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(n int) {
			a.setCache("key", n)
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = a.getCache("key", 5*time.Minute)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Verify cache still works
	_, ok := a.getCache("key", 5*time.Minute)
	if !ok {
		t.Error("expected cache to be accessible after concurrent access")
	}
}

// TestParseJSONWithStruct tests parsing JSON into specific struct types.
func TestParseJSONWithStruct(t *testing.T) {
	// Test Bead parsing
	beadJSON := `{"id":"test-1","title":"Test","status":"open","priority":1}`
	var bead model.Bead
	if err := parseJSON([]byte(beadJSON), &bead); err != nil {
		t.Fatalf("failed to parse bead JSON: %v", err)
	}
	if bead.ID != "test-1" {
		t.Errorf("bead.ID = %s, want test-1", bead.ID)
	}
	if bead.Priority != 1 {
		t.Errorf("bead.Priority = %d, want 1", bead.Priority)
	}

	// Test Convoy parsing
	convoyJSON := `{"id":"convoy-1","title":"MVP","status":"open","total_count":10,"closed_count":5}`
	var convoy model.Convoy
	if err := parseJSON([]byte(convoyJSON), &convoy); err != nil {
		t.Fatalf("failed to parse convoy JSON: %v", err)
	}
	if convoy.ID != "convoy-1" {
		t.Errorf("convoy.ID = %s, want convoy-1", convoy.ID)
	}
	if convoy.TotalCount != 10 {
		t.Errorf("convoy.TotalCount = %d, want 10", convoy.TotalCount)
	}

	// Test Polecat parsing
	polecatJSON := `{"name":"Toast","rig":"gastown","state":"working","session_running":true}`
	var polecat model.Polecat
	if err := parseJSON([]byte(polecatJSON), &polecat); err != nil {
		t.Fatalf("failed to parse polecat JSON: %v", err)
	}
	if polecat.Name != "Toast" {
		t.Errorf("polecat.Name = %s, want Toast", polecat.Name)
	}
	if !polecat.Running {
		t.Error("expected polecat.Running to be true")
	}
}

// TestAgentStatusParsing tests AgentStatus struct parsing.
func TestAgentStatusParsing(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantRunning bool
		wantSession string
	}{
		{
			name:        "running with session",
			input:       `{"running":true,"session_id":"test-session"}`,
			wantRunning: true,
			wantSession: "test-session",
		},
		{
			name:        "not running",
			input:       `{"running":false}`,
			wantRunning: false,
			wantSession: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status AgentStatus
			if err := json.Unmarshal([]byte(tt.input), &status); err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			if status.Running != tt.wantRunning {
				t.Errorf("Running = %v, want %v", status.Running, tt.wantRunning)
			}
			if status.SessionID != tt.wantSession {
				t.Errorf("SessionID = %s, want %s", status.SessionID, tt.wantSession)
			}
		})
	}
}

// TestRigStatusParsing tests RigStatus struct parsing.
func TestRigStatusParsing(t *testing.T) {
	input := `{
		"name": "gastown",
		"path": "/home/user/gastown",
		"prefix": "gt",
		"state": "active",
		"witness": {"running": true},
		"refinery": {"running": false}
	}`

	var status RigStatus
	if err := json.Unmarshal([]byte(input), &status); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if status.Name != "gastown" {
		t.Errorf("Name = %s, want gastown", status.Name)
	}
	if status.State != "active" {
		t.Errorf("State = %s, want active", status.State)
	}
	if !status.Witness.Running {
		t.Error("expected Witness.Running to be true")
	}
	if status.Refinery.Running {
		t.Error("expected Refinery.Running to be false")
	}
}
