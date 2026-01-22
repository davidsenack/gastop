package tui

import (
	"github.com/davidsenack/gastop/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// PolecatsPanel displays the polecats list.
type PolecatsPanel struct {
	list         *tview.List
	polecats     []model.Polecat
	selectedFunc func(*model.Polecat)
}

// NewPolecatsPanel creates a new polecats panel.
func NewPolecatsPanel() *PolecatsPanel {
	list := tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkBlue)

	list.SetBorder(true).SetTitle(" POLECATS ")

	return &PolecatsPanel{list: list}
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
		// Build primary text with state icon
		icon := pc.StateIcon()
		iconColor := ""
		if pc.Stuck {
			icon = "⚠"
			iconColor = "[red]"
		} else {
			switch pc.State {
			case "working":
				iconColor = "[blue]"
			case "done":
				iconColor = "[green]"
			case "idle":
				iconColor = "[gray]"
			}
		}

		primary := iconColor + icon + "[-] " + pc.Name
		if pc.Rig != "" {
			primary = iconColor + icon + "[-] " + pc.Rig + "/" + pc.Name
		}

		// Build secondary text with assigned bead and status
		secondary := "  " + pc.State
		if pc.AssignedBead != "" {
			secondary += " → " + pc.AssignedBead
		}
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
