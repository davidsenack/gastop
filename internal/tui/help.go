package tui

import (
	"github.com/rivo/tview"
)

// NewHelpOverlay creates a styled help overlay with keyboard shortcuts.
func NewHelpOverlay() *tview.TextView {
	helpText := `[::b]gastop - Gas Town Monitor[::-]
[darkgray]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[-]

[yellow::b]Navigation[::-]
  [aqua]j[white]/[aqua]k[-]             Move down/up in lists
  [aqua]h[white]/[aqua]l[-]             Switch panels left/right
  [aqua]g[white]/[aqua]G[-]             Jump to top/bottom of list
  [aqua]Tab[-]           Focus next panel
  [aqua]Shift-Tab[-]     Focus previous panel

[yellow::b]Actions[::-]
  [aqua]Enter[-]         Drill down / select item
  [aqua]x[white] or [aqua]d[-]         Kill polecat / close bead
  [aqua]r[-]             Manual refresh data
  [aqua]t[-]             Toggle auto-refresh on/off

[yellow::b]Display[::-]
  [aqua]L[-]             Toggle events/logs panel
  [aqua]+[white]/[aqua]=[-]           Faster refresh (min 1s)
  [aqua]-[-]             Slower refresh (max 30s)
  [aqua]/[-]             Search beads by ID or title
  [aqua]f[-]             Filter beads by status

[yellow::b]General[::-]
  [aqua]?[-]             Show this help
  [aqua]q[-]             Quit gastop

[darkgray]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[-]
[white::d]Press [aqua]Esc[white] or [aqua]Enter[white] to close[::-]`

	view := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText(helpText)

	view.SetBorder(true).
		SetTitle(" Keyboard Shortcuts ").
		SetTitleAlign(tview.AlignCenter)

	return view
}

// SearchModal wraps a form with its input field reference.
type SearchModal struct {
	*tview.Form
	inputField *tview.InputField
}

// NewSearchModal creates a search input modal.
func NewSearchModal(onSearch func(query string), onCancel func()) *SearchModal {
	form := tview.NewForm()
	modal := &SearchModal{Form: form}

	form.AddInputField("Search:", "", 40, nil, nil)
	modal.inputField = form.GetFormItemByLabel("Search:").(*tview.InputField)

	form.AddButton("Search", func() {
		query := modal.inputField.GetText()
		if onSearch != nil {
			onSearch(query)
		}
	})
	form.AddButton("Cancel", func() {
		if onCancel != nil {
			onCancel()
		}
	})

	form.SetBorder(true).SetTitle(" Search Beads (by ID or Title) ")
	form.SetCancelFunc(func() {
		if onCancel != nil {
			onCancel()
		}
	})

	return modal
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
