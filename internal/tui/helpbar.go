package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpBar displays keyboard shortcuts at the bottom of the screen.
type HelpBar struct {
	view *tview.TextView
}

// NewHelpBar creates a new help bar.
func NewHelpBar() *HelpBar {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	view.SetBackgroundColor(tcell.ColorDarkBlue)

	h := &HelpBar{view: view}
	h.UpdateDefault()
	return h
}

// Primitive returns the tview primitive.
func (h *HelpBar) Primitive() tview.Primitive {
	return h.view
}

// UpdateDefault shows the default keyboard shortcuts.
func (h *HelpBar) UpdateDefault() {
	shortcuts := "[::b]j/k[::-] Nav  [::b]h/l[::-] Panels  [::b]Tab[::-] Next  [::b]x[::-] Kill  [::b]r[::-] Refresh  [::b]+/-[::-] Speed  [::b]l[::-] Logs  [::b]?[::-] Help  [::b]q[::-] Quit"
	h.view.SetText(shortcuts)
}

// UpdateForPanel shows context-specific shortcuts for the given panel.
func (h *HelpBar) UpdateForPanel(panel string) {
	var shortcuts string
	switch panel {
	case "convoys":
		shortcuts = "[::b]j/k[::-] Navigate  [::b]Enter[::-] View beads  [::b]x[::-] Close convoy  [::b]h/l[::-] Switch panel  [::b]?[::-] Help  [::b]q[::-] Quit"
	case "beads":
		shortcuts = "[::b]j/k[::-] Navigate  [::b]Enter[::-] Details  [::b]x[::-] Close bead  [::b]h/l[::-] Switch panel  [::b]?[::-] Help  [::b]q[::-] Quit"
	case "polecats":
		shortcuts = "[::b]j/k[::-] Navigate  [::b]Enter[::-] Details  [::b]x[::-] Kill polecat  [::b]h/l[::-] Switch panel  [::b]?[::-] Help  [::b]q[::-] Quit"
	case "events":
		shortcuts = "[::b]j/k[::-] Scroll  [::b]h/l[::-] Switch panel  [::b]G[::-] Bottom  [::b]g[::-] Top  [::b]?[::-] Help  [::b]q[::-] Quit"
	default:
		h.UpdateDefault()
		return
	}
	h.view.SetText(shortcuts)
}
