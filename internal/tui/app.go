package tui

import (
	"context"
	"sync"
	"time"

	"github.com/davidsenack/gastop/internal/adapter"
	"github.com/davidsenack/gastop/internal/config"
	"github.com/davidsenack/gastop/internal/model"
	"github.com/davidsenack/gastop/internal/stuck"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the main beadtop application.
type App struct {
	app     *tview.Application
	adapter *adapter.Adapter
	config  *config.Config
	stuck   *stuck.Detector

	// Layout
	layout    *tview.Flex
	statusBar *StatusBar
	convoys   *ConvoysPanel
	beads     *BeadsPanel
	polecats  *PolecatsPanel
	events    *EventsPanel

	// State
	mu           sync.RWMutex
	convoyData   []model.Convoy
	beadData     []model.Bead
	polecatData  []model.Polecat
	eventData    []model.Event
	townStatus   *adapter.TownStatus
	currentRig   string
	autoRefresh  bool
	showLogs     bool
	lastError    string
	lastRefresh  time.Time

	// Context for background operations
	ctx    context.Context
	cancel context.CancelFunc
}

// NewApp creates a new beadtop application.
func NewApp(cfg *config.Config, adp *adapter.Adapter) *App {
	ctx, cancel := context.WithCancel(context.Background())

	a := &App{
		app:         tview.NewApplication(),
		adapter:     adp,
		config:      cfg,
		stuck:       stuck.NewDetector(cfg.StuckThresholdMins),
		autoRefresh: true,
		showLogs:    cfg.ShowLogs,
		ctx:         ctx,
		cancel:      cancel,
	}

	a.setupUI()
	a.setupKeyBindings()

	return a
}

// setupUI creates the UI layout.
func (a *App) setupUI() {
	// Create panels
	a.statusBar = NewStatusBar()
	a.convoys = NewConvoysPanel()
	a.beads = NewBeadsPanel()
	a.polecats = NewPolecatsPanel()
	a.events = NewEventsPanel(a.config.LogLines)

	// Wire up selection handlers
	a.convoys.SetSelectedFunc(func(convoy *model.Convoy) {
		// When a convoy is selected, filter beads to show its tracked issues
		if convoy != nil {
			a.filterBeadsByConvoy(convoy)
		}
	})

	// Create main content area (3 columns)
	mainContent := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.convoys.Primitive(), 0, 1, true).
		AddItem(a.beads.Primitive(), 0, 2, false).
		AddItem(a.polecats.Primitive(), 0, 1, false)

	// Main layout with optional logs panel
	a.layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.statusBar.Primitive(), 1, 0, false).
		AddItem(mainContent, 0, 1, true)

	if a.showLogs {
		a.layout.AddItem(a.events.Primitive(), a.config.LogLines+2, 0, false)
	}

	a.app.SetRoot(a.layout, true)
}

// setupKeyBindings configures global key bindings.
func (a *App) setupKeyBindings() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			a.cycleFocus()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				a.Stop()
				return nil
			case 'r':
				go a.refresh()
				return nil
			case 't':
				a.toggleAutoRefresh()
				return nil
			case 'l':
				a.toggleLogs()
				return nil
			case '?':
				a.showHelp()
				return nil
			case '/':
				a.showSearch()
				return nil
			case 'f':
				a.showFilter()
				return nil
			case '+', '=':
				a.decreaseRefreshInterval()
				return nil
			case '-':
				a.increaseRefreshInterval()
				return nil
			}
		}
		return event
	})
}

// cycleFocus cycles focus between panels.
func (a *App) cycleFocus() {
	focused := a.app.GetFocus()
	switch focused {
	case a.convoys.Primitive():
		a.app.SetFocus(a.beads.Primitive())
	case a.beads.Primitive():
		a.app.SetFocus(a.polecats.Primitive())
	case a.polecats.Primitive():
		if a.showLogs {
			a.app.SetFocus(a.events.Primitive())
		} else {
			a.app.SetFocus(a.convoys.Primitive())
		}
	case a.events.Primitive():
		a.app.SetFocus(a.convoys.Primitive())
	default:
		a.app.SetFocus(a.convoys.Primitive())
	}
}

// toggleAutoRefresh toggles automatic refresh.
func (a *App) toggleAutoRefresh() {
	a.mu.Lock()
	a.autoRefresh = !a.autoRefresh
	a.mu.Unlock()
	a.updateStatusBar()
}

// decreaseRefreshInterval makes refresh faster (min 1s).
func (a *App) decreaseRefreshInterval() {
	if a.config.RefreshInterval > time.Second {
		a.config.RefreshInterval -= time.Second
	}
	a.app.QueueUpdateDraw(func() {
		a.updateStatusBar()
	})
}

// increaseRefreshInterval makes refresh slower (max 30s).
func (a *App) increaseRefreshInterval() {
	if a.config.RefreshInterval < 30*time.Second {
		a.config.RefreshInterval += time.Second
	}
	a.app.QueueUpdateDraw(func() {
		a.updateStatusBar()
	})
}

// toggleLogs toggles the log panel visibility.
func (a *App) toggleLogs() {
	a.mu.Lock()
	a.showLogs = !a.showLogs
	a.mu.Unlock()

	// Rebuild layout
	a.app.QueueUpdateDraw(func() {
		mainContent := a.layout.GetItem(1)
		a.layout.Clear()
		a.layout.AddItem(a.statusBar.Primitive(), 1, 0, false)
		a.layout.AddItem(mainContent, 0, 1, true)
		if a.showLogs {
			a.layout.AddItem(a.events.Primitive(), a.config.LogLines+2, 0, false)
		}
	})
}

// showHelp displays the help modal.
func (a *App) showHelp() {
	help := NewHelpModal()
	help.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		a.app.SetRoot(a.layout, true)
	})
	a.app.SetRoot(help, true)
}

// showSearch displays the search dialog.
func (a *App) showSearch() {
	// TODO: Implement search modal
}

// showFilter displays the filter dialog.
func (a *App) showFilter() {
	// TODO: Implement filter modal
}

// filterBeadsByConvoy filters beads to show only those tracked by the convoy.
func (a *App) filterBeadsByConvoy(convoy *model.Convoy) {
	// For now, just update the title to show we're filtering
	// Full implementation would filter beadData
	a.beads.SetTitle("BEADS [" + convoy.ID + "]")
}

// refresh fetches new data from Gas Town.
// Updates UI incrementally as each data source returns.
func (a *App) refresh() {
	a.mu.Lock()
	a.lastRefresh = time.Now()
	a.mu.Unlock()

	// Clear error and update status bar immediately (spinner tick)
	a.mu.Lock()
	a.lastError = ""
	a.mu.Unlock()
	a.app.QueueUpdateDraw(func() {
		a.updateStatusBar()
	})

	// Fetch polecats (usually fastest)
	go func() {
		polecats, err := a.adapter.ListPolecats(a.ctx, "")
		if err != nil {
			return // Use cached data
		}
		a.stuck.CheckPolecats(polecats)
		a.mu.Lock()
		a.polecatData = polecats
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.polecats.Update(polecats)
		})
	}()

	// Fetch beads
	go func() {
		beads, err := a.adapter.ListBeads(a.ctx, adapter.BeadListOpts{Limit: 100})
		if err != nil {
			return // Use cached data
		}
		a.stuck.CheckBeads(beads)
		a.mu.Lock()
		a.beadData = beads
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.beads.Update(beads)
		})
	}()

	// Fetch convoys
	go func() {
		convoys, err := a.adapter.ListConvoys(a.ctx, adapter.ConvoyListOpts{})
		if err != nil {
			return // Use cached data
		}
		a.stuck.CheckConvoys(convoys, nil)
		a.mu.Lock()
		a.convoyData = convoys
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.convoys.Update(convoys)
		})
	}()

	// Fetch events (direct file read, very fast)
	go func() {
		events, err := a.adapter.TailEvents(a.ctx, a.config.LogLines)
		if err != nil {
			return // Events are optional
		}
		a.mu.Lock()
		a.eventData = events
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.events.Update(events)
		})
	}()

	// Skip gt status (too slow ~4s) - status bar updates from refresh tick above
	// Town name comes from config or is detected on startup
}

// updateStatusBar updates the status bar with current state.
func (a *App) updateStatusBar() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	townName := "Gas Town"
	if a.townStatus != nil && a.townStatus.Name != "" {
		townName = a.townStatus.Name
	}

	interval := a.config.RefreshInterval.String()
	if !a.autoRefresh {
		interval = "paused"
	}

	// Connected = we have some data
	connected := len(a.polecatData) > 0 || len(a.beadData) > 0 || a.lastError == ""

	a.statusBar.Update(townName, a.currentRig, interval, connected, false, a.lastError)
}

// Run starts the application.
func (a *App) Run() error {
	// Initial refresh
	go a.refresh()

	// Start refresh loop
	go a.refreshLoop()

	// Run the app
	return a.app.Run()
}

// refreshLoop periodically refreshes data.
func (a *App) refreshLoop() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-time.After(a.config.RefreshInterval):
			a.mu.RLock()
			autoRefresh := a.autoRefresh
			a.mu.RUnlock()

			if autoRefresh {
				a.refresh()
			}
		}
	}
}

// Stop stops the application.
func (a *App) Stop() {
	a.cancel()
	a.app.Stop()
}
