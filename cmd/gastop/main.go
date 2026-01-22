package main

import (
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
	// TODO: Implement JSON output mode
	fmt.Println(`{"status": "not implemented"}`)
}
