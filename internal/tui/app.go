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

// App is the main gastop application.
type App struct {
	app     *tview.Application
	adapter *adapter.Adapter
	config  *config.Config
	stuck   *stuck.Detector

	// Layout
	layout      *tview.Flex
	mainContent *tview.Flex
	statusBar   *StatusBar
	helpBar     *HelpBar
	convoys     *ConvoysPanel
	beads       *BeadsPanel
	polecats    *PolecatsPanel
	events      *EventsPanel

	// Panel tracking for vim navigation
	panels       []tview.Primitive
	currentPanel int

	// State
	mu               sync.RWMutex
	convoyData       []model.Convoy
	beadData         []model.Bead
	polecatData      []model.Polecat
	eventData        []model.Event
	townStatus       *adapter.TownStatus
	currentRig       string
	autoRefresh      bool
	showLogs         bool
	lastError        string
	lastRefresh      time.Time
	beadStatusFilter string // Filter beads by status ("" = all)

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
	a.helpBar = NewHelpBar()
	a.convoys = NewConvoysPanel()
	a.beads = NewBeadsPanel()
	a.polecats = NewPolecatsPanel()
	a.events = NewEventsPanel(a.config.LogLines)

	// Track panels for vim navigation (h/l)
	a.panels = []tview.Primitive{
		a.convoys.Primitive(),
		a.beads.Primitive(),
		a.polecats.Primitive(),
	}
	if a.showLogs {
		a.panels = append(a.panels, a.events.Primitive())
	}
	a.currentPanel = 0

	// Wire up selection handlers
	a.convoys.SetSelectedFunc(func(convoy *model.Convoy) {
		// When a convoy is selected, filter beads to show its tracked issues
		if convoy != nil {
			a.filterBeadsByConvoy(convoy)
		}
	})

	// Create main content area (3 columns)
	a.mainContent = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.convoys.Primitive(), 0, 1, true).
		AddItem(a.beads.Primitive(), 0, 2, false).
		AddItem(a.polecats.Primitive(), 0, 1, false)

	// Main layout with help bar at bottom
	a.layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.statusBar.Primitive(), 1, 0, false).
		AddItem(a.mainContent, 0, 1, true)

	if a.showLogs {
		a.layout.AddItem(a.events.Primitive(), a.config.LogLines+2, 0, false)
	}

	// Add help bar at the very bottom
	a.layout.AddItem(a.helpBar.Primitive(), 1, 0, false)

	a.app.SetRoot(a.layout, true)
}

// setupKeyBindings configures global key bindings.
func (a *App) setupKeyBindings() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			a.focusNext()
			return nil
		case tcell.KeyBacktab:
			a.focusPrev()
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
			case 'L': // Capital L for logs toggle (lowercase l is vim right)
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
			// Vim-style panel navigation
			case 'h':
				a.focusPrev()
				return nil
			case 'l':
				a.focusNext()
				return nil
			// Vim-style list navigation (j/k handled by list itself, but we ensure it works)
			case 'j':
				a.navigateDown()
				return nil
			case 'k':
				a.navigateUp()
				return nil
			case 'g':
				a.navigateTop()
				return nil
			case 'G':
				a.navigateBottom()
				return nil
			// Kill/close action
			case 'x', 'd':
				a.killSelected()
				return nil
			}
		}
		return event
	})
}

// focusNext moves focus to the next panel (vim l / Tab).
func (a *App) focusNext() {
	a.currentPanel = (a.currentPanel + 1) % len(a.panels)
	a.app.SetFocus(a.panels[a.currentPanel])
	a.updateHelpBarForFocus()
}

// focusPrev moves focus to the previous panel (vim h / Shift-Tab).
func (a *App) focusPrev() {
	a.currentPanel--
	if a.currentPanel < 0 {
		a.currentPanel = len(a.panels) - 1
	}
	a.app.SetFocus(a.panels[a.currentPanel])
	a.updateHelpBarForFocus()
}

// updateHelpBarForFocus updates help bar based on current focus.
func (a *App) updateHelpBarForFocus() {
	focused := a.app.GetFocus()
	switch focused {
	case a.convoys.Primitive():
		a.helpBar.UpdateForPanel("convoys")
	case a.beads.Primitive():
		a.helpBar.UpdateForPanel("beads")
	case a.polecats.Primitive():
		a.helpBar.UpdateForPanel("polecats")
	case a.events.Primitive():
		a.helpBar.UpdateForPanel("events")
	default:
		a.helpBar.UpdateDefault()
	}
}

// navigateDown moves selection down in current list (vim j).
func (a *App) navigateDown() {
	focused := a.app.GetFocus()
	switch focused {
	case a.convoys.Primitive():
		a.convoys.list.SetCurrentItem(a.convoys.list.GetCurrentItem() + 1)
	case a.beads.Primitive():
		row, _ := a.beads.table.GetSelection()
		if row < a.beads.table.GetRowCount()-1 {
			a.beads.table.Select(row+1, 0)
		}
	case a.polecats.Primitive():
		a.polecats.list.SetCurrentItem(a.polecats.list.GetCurrentItem() + 1)
	case a.events.Primitive():
		a.events.ScrollDown()
	}
}

// navigateUp moves selection up in current list (vim k).
func (a *App) navigateUp() {
	focused := a.app.GetFocus()
	switch focused {
	case a.convoys.Primitive():
		idx := a.convoys.list.GetCurrentItem() - 1
		if idx >= 0 {
			a.convoys.list.SetCurrentItem(idx)
		}
	case a.beads.Primitive():
		row, _ := a.beads.table.GetSelection()
		if row > 1 { // Skip header row
			a.beads.table.Select(row-1, 0)
		}
	case a.polecats.Primitive():
		idx := a.polecats.list.GetCurrentItem() - 1
		if idx >= 0 {
			a.polecats.list.SetCurrentItem(idx)
		}
	case a.events.Primitive():
		a.events.ScrollUp()
	}
}

// navigateTop moves to top of list (vim g).
func (a *App) navigateTop() {
	focused := a.app.GetFocus()
	switch focused {
	case a.convoys.Primitive():
		a.convoys.list.SetCurrentItem(0)
	case a.beads.Primitive():
		a.beads.table.Select(1, 0) // Skip header
	case a.polecats.Primitive():
		a.polecats.list.SetCurrentItem(0)
	case a.events.Primitive():
		a.events.ScrollToTop()
	}
}

// navigateBottom moves to bottom of list (vim G).
func (a *App) navigateBottom() {
	focused := a.app.GetFocus()
	switch focused {
	case a.convoys.Primitive():
		a.convoys.list.SetCurrentItem(a.convoys.list.GetItemCount() - 1)
	case a.beads.Primitive():
		a.beads.table.Select(a.beads.table.GetRowCount()-1, 0)
	case a.polecats.Primitive():
		a.polecats.list.SetCurrentItem(a.polecats.list.GetItemCount() - 1)
	case a.events.Primitive():
		a.events.ScrollToBottom()
	}
}

// killSelected kills/closes the selected item based on current panel.
func (a *App) killSelected() {
	focused := a.app.GetFocus()
	switch focused {
	case a.polecats.Primitive():
		if pc := a.polecats.Selected(); pc != nil {
			a.showConfirmKillPolecat(pc)
		}
	case a.beads.Primitive():
		if b := a.beads.Selected(); b != nil {
			a.showConfirmCloseBead(b)
		}
	case a.convoys.Primitive():
		// Convoys can't be directly killed, show message
		a.showMessage("Convoys close automatically when all beads complete")
	}
}

// showConfirmKillPolecat shows a confirmation dialog for killing a polecat.
func (a *App) showConfirmKillPolecat(pc *model.Polecat) {
	name := pc.Name
	if pc.Rig != "" {
		name = pc.Rig + "/" + pc.Name
	}

	modal := tview.NewModal().
		SetText("Kill polecat " + name + "?\n\nThis will terminate the session and remove the worktree.").
		AddButtons([]string{"Cancel", "Kill"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Kill" {
				go func() {
					err := a.adapter.NukePolecat(a.ctx, pc.Rig, pc.Name)
					a.app.QueueUpdateDraw(func() {
						if err != nil {
							a.showMessage("Failed to kill polecat: " + err.Error())
						} else {
							a.showMessage("Polecat " + name + " killed")
							go a.refresh()
						}
					})
				}()
			}
			a.app.SetRoot(a.layout, true)
		})
	a.app.SetRoot(modal, true)
}

// showConfirmCloseBead shows a confirmation dialog for closing a bead.
func (a *App) showConfirmCloseBead(b *model.Bead) {
	modal := tview.NewModal().
		SetText("Close bead " + b.ID + "?\n\n" + b.Title).
		AddButtons([]string{"Cancel", "Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Close" {
				go func() {
					err := a.adapter.CloseBead(a.ctx, b.ID)
					a.app.QueueUpdateDraw(func() {
						if err != nil {
							a.showMessage("Failed to close bead: " + err.Error())
						} else {
							a.showMessage("Bead " + b.ID + " closed")
							go a.refresh()
						}
					})
				}()
			}
			a.app.SetRoot(a.layout, true)
		})
	a.app.SetRoot(modal, true)
}

// showMessage shows a temporary message modal.
func (a *App) showMessage(msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.app.SetRoot(a.layout, true)
		})
	a.app.SetRoot(modal, true)
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
	showLogs := a.showLogs
	a.mu.Unlock()

	// Update panels list
	a.panels = []tview.Primitive{
		a.convoys.Primitive(),
		a.beads.Primitive(),
		a.polecats.Primitive(),
	}
	if showLogs {
		a.panels = append(a.panels, a.events.Primitive())
	}

	// Rebuild layout
	a.app.QueueUpdateDraw(func() {
		a.layout.Clear()
		a.layout.AddItem(a.statusBar.Primitive(), 1, 0, false)
		a.layout.AddItem(a.mainContent, 0, 1, true)
		if showLogs {
			a.layout.AddItem(a.events.Primitive(), a.config.LogLines+2, 0, false)
		}
		a.layout.AddItem(a.helpBar.Primitive(), 1, 0, false)
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
	searchModal := NewSearchModal(
		func(query string) {
			// Perform search and return to main view
			a.beads.Search(query)
			a.app.SetRoot(a.layout, true)
			a.app.SetFocus(a.beads.Primitive())
		},
		func() {
			// Cancel - return to main view
			a.app.SetRoot(a.layout, true)
		},
	)

	// Center the search form in a flex container
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(searchModal, 50, 0, true).
			AddItem(nil, 0, 1, false), 7, 0, true).
		AddItem(nil, 0, 1, false)

	a.app.SetRoot(flex, true)
}

// showFilter displays the filter dialog.
func (a *App) showFilter() {
	// Create filter list
	list := tview.NewList().
		AddItem("All", "Show all beads", 'a', func() {
			a.applyBeadFilter("")
		}).
		AddItem("Open", "Show open beads", 'o', func() {
			a.applyBeadFilter("open")
		}).
		AddItem("In Progress", "Show in_progress beads", 'i', func() {
			a.applyBeadFilter("in_progress")
		}).
		AddItem("Blocked", "Show blocked beads", 'b', func() {
			a.applyBeadFilter("blocked")
		}).
		AddItem("Deferred", "Show deferred beads", 'd', func() {
			a.applyBeadFilter("deferred")
		}).
		AddItem("Closed", "Show closed beads", 'c', func() {
			a.applyBeadFilter("closed")
		})

	list.SetBorder(true).SetTitle(" Filter by Status ")

	// Handle escape key to close
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.app.SetRoot(a.layout, true)
			return nil
		}
		return event
	})

	// Center the list in a flex container for modal-like appearance
	modal := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(list, 40, 0, true).
			AddItem(nil, 0, 1, false), 12, 0, true).
		AddItem(nil, 0, 1, false)

	a.app.SetRoot(modal, true)
}

// applyBeadFilter applies a status filter to beads and returns to main layout.
func (a *App) applyBeadFilter(status string) {
	a.mu.Lock()
	a.beadStatusFilter = status
	beads := a.beadData
	a.mu.Unlock()

	// Filter beads
	filtered := a.filterBeadsByStatus(beads, status)

	// Update panel
	a.app.QueueUpdateDraw(func() {
		a.beads.Update(filtered)
		if status == "" {
			a.beads.SetTitle("BEADS")
		} else {
			a.beads.SetTitle("BEADS [" + status + "]")
		}
		a.app.SetRoot(a.layout, true)
	})
}

// filterBeadsByStatus filters a bead slice by status.
func (a *App) filterBeadsByStatus(beads []model.Bead, status string) []model.Bead {
	if status == "" {
		return beads
	}
	var filtered []model.Bead
	for _, b := range beads {
		if b.Status == status {
			filtered = append(filtered, b)
		}
	}
	return filtered
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
			a.mu.Lock()
			a.lastError = "polecats: " + err.Error()
			a.mu.Unlock()
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
		filter := a.beadStatusFilter
		a.mu.Unlock()

		// Apply current filter
		filtered := a.filterBeadsByStatus(beads, filter)

		a.app.QueueUpdateDraw(func() {
			a.beads.Update(filtered)
			if filter != "" {
				a.beads.SetTitle("BEADS [" + filter + "]")
			}
		})
	}()

	// Fetch convoys
	go func() {
		convoys, err := a.adapter.ListConvoys(a.ctx, adapter.ConvoyListOpts{})
		if err != nil {
			a.mu.Lock()
			a.lastError = "convoys: " + err.Error()
			a.mu.Unlock()
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
