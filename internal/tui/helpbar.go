package tui

import (
	"github.com/rivo/tview"
)

// HelpBar displays keyboard shortcuts at the bottom of the screen.
type HelpBar struct {
	view *tview.TextView
}

// NewHelpBar creates a new help bar.
func NewHelpBar() *HelpBar {
	theme := GetTheme()
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(theme.Foreground)

	view.SetBackgroundColor(theme.SelectionBg)

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
	tags := GetTags()
	key := "[" + tags.Accent1 + "][::b]"
	end := "[::-][-]"
	shortcuts := key + "j/k" + end + " Nav  " + key + "h/l" + end + " Panels  " + key + "Tab" + end + " Next  " + key + "x" + end + " Kill  " + key + "r" + end + " Refresh  " + key + "+" + end + "/" + key + "-" + end + " Speed  " + key + "/" + end + " Search  " + key + "f" + end + " Filter  " + key + "?" + end + " Help  " + key + "q" + end + " Quit"
	h.view.SetText(shortcuts)
}

// UpdateForPanel shows context-specific shortcuts for the given panel.
func (h *HelpBar) UpdateForPanel(panel string) {
	tags := GetTags()
	key := "[" + tags.Accent1 + "][::b]"
	end := "[::-][-]"

	var shortcuts string
	switch panel {
	case "convoys":
		shortcuts = key + "j/k" + end + " Navigate  " + key + "Enter" + end + " View beads  " + key + "x" + end + " Close convoy  " + key + "h/l" + end + " Switch panel  " + key + "?" + end + " Help  " + key + "q" + end + " Quit"
	case "beads":
		shortcuts = key + "j/k" + end + " Navigate  " + key + "Enter" + end + " Details  " + key + "x" + end + " Close bead  " + key + "/" + end + " Search  " + key + "f" + end + " Filter  " + key + "?" + end + " Help  " + key + "q" + end + " Quit"
	case "polecats":
		shortcuts = key + "j/k" + end + " Navigate  " + key + "Enter" + end + " Details  " + key + "x" + end + " Kill polecat  " + key + "h/l" + end + " Switch panel  " + key + "?" + end + " Help  " + key + "q" + end + " Quit"
	case "events":
		shortcuts = key + "j/k" + end + " Scroll  " + key + "h/l" + end + " Switch panel  " + key + "G" + end + " Bottom  " + key + "g" + end + " Top  " + key + "?" + end + " Help  " + key + "q" + end + " Quit"
	default:
		h.UpdateDefault()
		return
	}
	h.view.SetText(shortcuts)
}
