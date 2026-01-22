package tui

import (
	"fmt"

	"github.com/davidsenack/gastop/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Spinner frames for working animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// PolecatsPanel displays the polecats list.
type PolecatsPanel struct {
	list         *tview.List
	polecats     []model.Polecat
	selectedFunc func(*model.Polecat)
	spinnerIdx   int
}

// NewPolecatsPanel creates a new polecats panel.
func NewPolecatsPanel() *PolecatsPanel {
	list := tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkBlue)

	list.SetBorder(true).SetTitle(" POLECATS ")

	return &PolecatsPanel{list: list, spinnerIdx: 0}
}

// AdvanceSpinner advances the spinner animation frame.
func (p *PolecatsPanel) AdvanceSpinner() {
	p.spinnerIdx = (p.spinnerIdx + 1) % len(spinnerFrames)
}

// Primitive returns the tview primitive.
func (p *PolecatsPanel) Primitive() tview.Primitive {
	return p.list
}

// Update updates the panel with new polecat data.
func (p *PolecatsPanel) Update(polecats []model.Polecat) {
	p.polecats = polecats
	currentIndex := p.list.GetCurrentItem()

	p.list.Clear()
	for i, pc := range polecats {
		// Build primary text with state icon and name
		icon := pc.StateIcon()
		iconColor := ""
		if pc.Stuck {
			icon = "⚠"
			iconColor = "[red]"
		} else {
			switch pc.State {
			case "working":
				// Use spinner for working state
				icon = spinnerFrames[p.spinnerIdx]
				iconColor = "[yellow]"
			case "done":
				iconColor = "[green]"
			case "idle":
				iconColor = "[gray]"
			}
		}

		// Primary line: icon + rig/name + session indicator
		primary := iconColor + icon + "[-] "
		if pc.Rig != "" {
			primary += pc.Rig + "/" + pc.Name
		} else {
			primary += pc.Name
		}

		// Add session status indicator
		if !pc.Running {
			primary += " [red]●[-]" // Red dot = stopped
		} else if pc.Attached {
			primary += " [blue]◉[-]" // Blue ring = attached
		}

		// Add activity time if available
		if activityAgo := pc.ActivityAgo(); activityAgo != "" {
			primary += " [gray](" + activityAgo + ")[-]"
		}

		// Build secondary text with work info
		secondary := "  "

		// Show work description (hooked bead, assigned bead, or branch)
		if work := pc.WorkDescription(); work != "" {
			secondary += "[white]" + work + "[-]"
		} else {
			secondary += "[gray]" + pc.State + "[-]"
		}

		// Add stuck reason if stuck
		if pc.Stuck {
			secondary += " [red]" + pc.StuckReason + "[-]"
		}

		idx := i
		p.list.AddItem(primary, secondary, 0, func() {
			if p.selectedFunc != nil && idx < len(p.polecats) {
				p.selectedFunc(&p.polecats[idx])
			}
		})
	}

	// Show placeholder if no polecats
	if len(polecats) == 0 {
		p.list.AddItem("[gray]No active polecats[-]", "", 0, nil)
	}

	// Restore selection
	if currentIndex >= 0 && currentIndex < len(polecats) {
		p.list.SetCurrentItem(currentIndex)
	}
}

// UpdateWithSpinner updates the panel and advances the spinner.
func (p *PolecatsPanel) UpdateWithSpinner(polecats []model.Polecat) {
	p.AdvanceSpinner()
	p.Update(polecats)
}

// CountByState returns a summary string of polecats by state.
func (p *PolecatsPanel) CountByState() string {
	working, done, idle := 0, 0, 0
	for _, pc := range p.polecats {
		switch pc.State {
		case "working":
			working++
		case "done":
			done++
		case "idle":
			idle++
		}
	}
	return fmt.Sprintf("%d working, %d done, %d idle", working, done, idle)
}

// SetSelectedFunc sets the callback for when a polecat is selected.
func (p *PolecatsPanel) SetSelectedFunc(fn func(*model.Polecat)) {
	p.selectedFunc = fn
}

// Selected returns the currently selected polecat.
func (p *PolecatsPanel) Selected() *model.Polecat {
	idx := p.list.GetCurrentItem()
	if idx >= 0 && idx < len(p.polecats) {
		return &p.polecats[idx]
	}
	return nil
}

// SetTitle sets the panel title.
func (p *PolecatsPanel) SetTitle(title string) {
	p.list.SetTitle(" " + title + " ")
}
