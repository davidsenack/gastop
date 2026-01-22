package tui

import (
	"github.com/rivo/tview"
)

// NewHelpModal creates a help modal with keyboard shortcuts.
func NewHelpModal() *tview.Modal {
	helpText := `gastop - Gas Town Monitor

Vim Navigation:
  j/k           Move down/up in lists
  h/l           Switch panels left/right
  g/G           Jump to top/bottom
  Tab/Shift-Tab Next/previous panel

Actions:
  x or d        Kill polecat / close bead
  Enter         Drill down / select
  r             Manual refresh
  t             Toggle auto-refresh

Display:
  L             Toggle logs panel
  +/=           Faster refresh (min 1s)
  -             Slower refresh (max 30s)
  /             Search
  f             Filter by status

General:
  ?             Show this help
  q             Quit

Press any key to close`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"})

	return modal
}

// NewSearchModal creates a search input modal.
func NewSearchModal(onSearch func(query string)) *tview.Form {
	form := tview.NewForm().
		AddInputField("Search:", "", 40, nil, nil).
		AddButton("Search", func() {
			// Get the search query and call callback
		}).
		AddButton("Cancel", nil)

	form.SetBorder(true).SetTitle(" Search ")

	return form
}

// NewFilterModal creates a filter selection modal.
func NewFilterModal(onFilter func(status string)) *tview.List {
	list := tview.NewList().
		AddItem("All", "Show all beads", 'a', func() { onFilter("") }).
		AddItem("Open", "Show open beads", 'o', func() { onFilter("open") }).
		AddItem("In Progress", "Show in_progress beads", 'i', func() { onFilter("in_progress") }).
		AddItem("Blocked", "Show blocked beads", 'b', func() { onFilter("blocked") }).
		AddItem("Closed", "Show closed beads", 'c', func() { onFilter("closed") })

	list.SetBorder(true).SetTitle(" Filter by Status ")

	return list
}
