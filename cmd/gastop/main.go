package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/davidsenack/gastop/internal/adapter"
	"github.com/davidsenack/gastop/internal/config"
	"github.com/davidsenack/gastop/internal/tui"
	flag "github.com/spf13/pflag"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	var (
		townRoot    = flag.StringP("town", "t", "", "Gas Town root directory (auto-detect if empty)")
		rig         = flag.StringP("rig", "r", "", "Focus on a specific rig")
		showVersion = flag.BoolP("version", "V", false, "Show version information")
		jsonOutput  = flag.BoolP("json", "j", false, "Output JSON instead of TUI (for scripting)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `gastop - htop-like TUI for Gas Town

Usage:
  gastop [flags]

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  gastop                       # Launch TUI, auto-detect town
  gastop -t ~/gt               # Specify town root (short)
  gastop --town ~/gt           # Specify town root (long)
  gastop -r gastown            # Focus on specific rig
  gastop -j                    # JSON output for scripting

Keyboard:
  j/k     Navigate up/down
  h/l     Switch panels
  x       Kill/close selected
  r       Refresh
  /       Search
  f       Filter
  ?       Help
  q       Quit
`)
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("gastop %s (%s)\n", version, commit)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Override config with flags
	if *townRoot != "" {
		cfg.Paths.TownRoot = *townRoot
	}

	// Create adapter
	adp := adapter.New(cfg.Paths.GTBinary, cfg.Paths.BDBinary, cfg.Paths.TownRoot)

	if *jsonOutput {
		// JSON mode - just dump data and exit
		runJSONMode(adp, *rig)
		return
	}

	// Create and run TUI
	app := tui.NewApp(cfg, adp)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runJSONMode(adp *adapter.Adapter, rig string) {
	ctx := context.Background()

	// Build JSON output structure
	output := struct {
		Status   *adapter.TownStatus `json:"status,omitempty"`
		Polecats interface{}         `json:"polecats,omitempty"`
		Beads    interface{}         `json:"beads,omitempty"`
		Convoys  interface{}         `json:"convoys,omitempty"`
		Error    string              `json:"error,omitempty"`
	}{}

	// Get town status
	status, err := adp.GetTownStatus(ctx)
	if err != nil {
		output.Error = fmt.Sprintf("failed to get town status: %v", err)
	} else {
		output.Status = status
	}

	// Get polecats
	polecats, err := adp.ListPolecats(ctx, rig)
	if err != nil {
		if output.Error != "" {
			output.Error += "; "
		}
		output.Error += fmt.Sprintf("failed to get polecats: %v", err)
	} else {
		output.Polecats = polecats
	}

	// Get beads
	beads, err := adp.ListBeads(ctx, adapter.BeadListOpts{})
	if err != nil {
		if output.Error != "" {
			output.Error += "; "
		}
		output.Error += fmt.Sprintf("failed to get beads: %v", err)
	} else {
		output.Beads = beads
	}

	// Get convoys
	convoys, err := adp.ListConvoys(ctx, adapter.ConvoyListOpts{})
	if err != nil {
		if output.Error != "" {
			output.Error += "; "
		}
		output.Error += fmt.Sprintf("failed to get convoys: %v", err)
	} else {
		output.Convoys = convoys
	}

	// Marshal and print JSON
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}
