package tui

import (
	"fmt"

	"github.com/davidsenack/gastop/internal/model"
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
	theme := GetTheme()
	list := tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(theme.SelectionBg).
		SetSelectedTextColor(theme.SelectionFg).
		SetMainTextColor(theme.Foreground).
		SetSecondaryTextColor(theme.Muted)

	list.SetBorder(true).
		SetTitle(" CONVOYS ").
		SetBorderColor(theme.BorderColor).
		SetTitleColor(theme.TitleColor)

	return &ConvoysPanel{list: list}
}

// Primitive returns the tview primitive.
func (p *ConvoysPanel) Primitive() tview.Primitive {
	return p.list
}

// Update updates the panel with new convoy data.
func (p *ConvoysPanel) Update(convoys []model.Convoy) {
	tags := GetTags()
	p.convoys = convoys
	currentIndex := p.list.GetCurrentItem()

	p.list.Clear()
	for i, c := range convoys {
		// Build primary text with status icon
		var icon string
		if c.Stuck {
			icon = "[" + tags.Error + "]⚠[-]"
		} else if c.Status == "closed" {
			icon = "[" + tags.Done + "]✓[-]"
		} else {
			icon = "[" + tags.Accent1 + "]●[-]"
		}

		primary := fmt.Sprintf("%s %s", icon, c.Title)
		if len(primary) > 25 {
			primary = primary[:22] + "..."
		}

		// Build secondary text with progress
		secondary := fmt.Sprintf("  [%s]%s[-] ", tags.Dim, c.ID)
		if c.TotalCount > 0 {
			// Show progress bar
			pct := float64(c.ClosedCount) / float64(c.TotalCount)
			secondary += renderProgressBar(pct, 8) + " "
		}
		secondary += fmt.Sprintf("[%s]%d/%d[-]", tags.Muted, c.ClosedCount, c.TotalCount)
		if c.Stuck {
			secondary += " [" + tags.Error + "]STUCK[-]"
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

// renderProgressBar renders a simple progress bar.
func renderProgressBar(pct float64, width int) string {
	tags := GetTags()
	filled := int(pct * float64(width))
	if filled > width {
		filled = width
	}

	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "[" + tags.Success + "]█[-]"
		} else {
			bar += "[" + tags.Dim + "]░[-]"
		}
	}
	bar += "]"
	return bar
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
