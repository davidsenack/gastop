package tui

import (
	"fmt"

	"github.com/davidsenack/gastop/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ConvoysPanel displays the convoys list.
type ConvoysPanel struct {
	list         *tview.List
	convoys      []model.Convoy
	selectedFunc func(*model.Convoy)
}

// NewConvoysPanel creates a new convoys panel.
func NewConvoysPanel() *ConvoysPanel {
	list := tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkBlue)

	list.SetBorder(true).SetTitle(" CONVOYS ")

	return &ConvoysPanel{list: list}
}

// Primitive returns the tview primitive.
func (p *ConvoysPanel) Primitive() tview.Primitive {
	return p.list
}

// Update updates the panel with new convoy data.
func (p *ConvoysPanel) Update(convoys []model.Convoy) {
	p.convoys = convoys
	currentIndex := p.list.GetCurrentItem()

	p.list.Clear()
	for i, c := range convoys {
		// Build primary text with status icon
		icon := c.StatusIcon()
		if c.Stuck {
			icon = "[red]⚠[-]"
		} else if c.Status == "closed" {
			icon = "[green]✓[-]"
		} else {
			icon = "[blue]●[-]"
		}

		primary := fmt.Sprintf("%s %s", icon, c.Title)
		if len(primary) > 25 {
			primary = primary[:22] + "..."
		}

		// Build secondary text with progress
		secondary := fmt.Sprintf("  %s [%d/%d]", c.ID, c.ClosedCount, c.TotalCount)
		if c.Stuck {
			secondary += " [red]STUCK[-]"
		}

		idx := i
		p.list.AddItem(primary, secondary, 0, func() {
			if p.selectedFunc != nil && idx < len(p.convoys) {
				p.selectedFunc(&p.convoys[idx])
			}
		})
	}

	// Restore selection
	if currentIndex >= 0 && currentIndex < len(convoys) {
		p.list.SetCurrentItem(currentIndex)
	}
}

// SetSelectedFunc sets the callback for when a convoy is selected.
func (p *ConvoysPanel) SetSelectedFunc(fn func(*model.Convoy)) {
	p.selectedFunc = fn
}

// Selected returns the currently selected convoy.
func (p *ConvoysPanel) Selected() *model.Convoy {
	idx := p.list.GetCurrentItem()
	if idx >= 0 && idx < len(p.convoys) {
		return &p.convoys[idx]
	}
	return nil
}

// SetTitle sets the panel title.
func (p *ConvoysPanel) SetTitle(title string) {
	p.list.SetTitle(" " + title + " ")
}
