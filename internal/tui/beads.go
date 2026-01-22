package tui

import (
	"fmt"
	"strings"

	"github.com/davidsenack/gastop/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BeadsPanel displays the beads table.
type BeadsPanel struct {
	table        *tview.Table
	beads        []model.Bead
	allBeads     []model.Bead // Original unfiltered data
	selectedFunc func(*model.Bead)
}

// NewBeadsPanel creates a new beads panel.
func NewBeadsPanel() *BeadsPanel {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkBlue))

	table.SetBorder(true).SetTitle(" BEADS ")

	// Set up header row
	headers := []string{"", "ID", "Status", "Pri", "Title", "Age"}
	for i, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(0)
		if i == 4 { // Title column expands
			cell.SetExpansion(1)
		}
		table.SetCell(0, i, cell)
	}

	return &BeadsPanel{table: table}
}

// Primitive returns the tview primitive.
func (p *BeadsPanel) Primitive() tview.Primitive {
	return p.table
}

// Update updates the panel with new bead data.
func (p *BeadsPanel) Update(beads []model.Bead) {
	p.allBeads = beads
	p.updateDisplay(beads)
}

// SetSelectedFunc sets the callback for when a bead is selected.
func (p *BeadsPanel) SetSelectedFunc(fn func(*model.Bead)) {
	p.selectedFunc = fn
	p.table.SetSelectedFunc(func(row, col int) {
		if row > 0 && row-1 < len(p.beads) && fn != nil {
			fn(&p.beads[row-1])
		}
	})
}

// Selected returns the currently selected bead.
func (p *BeadsPanel) Selected() *model.Bead {
	row, _ := p.table.GetSelection()
	if row > 0 && row-1 < len(p.beads) {
		return &p.beads[row-1]
	}
	return nil
}

// SetTitle sets the panel title.
func (p *BeadsPanel) SetTitle(title string) {
	p.table.SetTitle(" " + title + " ")
}

// FilterByStatus filters the displayed beads by status.
func (p *BeadsPanel) FilterByStatus(status string) {
	// Store original beads and filter
	var filtered []model.Bead
	for _, b := range p.allBeads {
		if status == "" || b.Status == status {
			filtered = append(filtered, b)
		}
	}
	p.updateDisplay(filtered)
	if status == "" {
		p.SetTitle("BEADS")
	} else {
		p.SetTitle(fmt.Sprintf("BEADS [%s]", status))
	}
}

// Search filters beads by ID or Title (case-insensitive).
func (p *BeadsPanel) Search(query string) {
	if query == "" {
		p.updateDisplay(p.allBeads)
		p.SetTitle("BEADS")
		return
	}

	query = strings.ToLower(query)
	var filtered []model.Bead
	for _, b := range p.allBeads {
		if strings.Contains(strings.ToLower(b.ID), query) ||
			strings.Contains(strings.ToLower(b.Title), query) {
			filtered = append(filtered, b)
		}
	}
	p.updateDisplay(filtered)
	p.SetTitle(fmt.Sprintf("BEADS [search: %s]", query))
}

// ClearSearch resets the beads display to show all beads.
func (p *BeadsPanel) ClearSearch() {
	p.updateDisplay(p.allBeads)
	p.SetTitle("BEADS")
}

// updateDisplay updates the table without changing allBeads.
func (p *BeadsPanel) updateDisplay(beads []model.Bead) {
	p.beads = beads

	// Remember current selection
	currentRow, _ := p.table.GetSelection()

	// Clear all rows except header
	for i := p.table.GetRowCount() - 1; i > 0; i-- {
		p.table.RemoveRow(i)
	}

	// Add bead rows
	for i, b := range beads {
		row := i + 1 // Skip header

		// Status icon
		icon := b.StatusIcon()
		iconColor := tcell.ColorWhite
		if b.Stuck {
			icon = "⚠"
			iconColor = tcell.ColorRed
		} else {
			switch b.Status {
			case "in_progress":
				iconColor = tcell.ColorBlue
			case "closed":
				iconColor = tcell.ColorGreen
			case "blocked":
				iconColor = tcell.ColorRed
			case "deferred":
				iconColor = tcell.ColorYellow
			}
		}

		p.table.SetCell(row, 0, tview.NewTableCell(icon).SetTextColor(iconColor))
		p.table.SetCell(row, 1, tview.NewTableCell(b.ID).SetTextColor(tcell.ColorDarkCyan))
		p.table.SetCell(row, 2, tview.NewTableCell(b.Status))
		p.table.SetCell(row, 3, tview.NewTableCell(b.PriorityString()))

		// Truncate title if needed
		title := b.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		titleCell := tview.NewTableCell(title).SetExpansion(1)
		if b.Stuck {
			titleCell.SetTextColor(tcell.ColorRed)
		}
		p.table.SetCell(row, 4, titleCell)

		p.table.SetCell(row, 5, tview.NewTableCell(b.Age))
	}

	// Restore selection
	if currentRow > 0 && currentRow <= len(beads) {
		p.table.Select(currentRow, 0)
	} else if len(beads) > 0 {
		p.table.Select(1, 0)
	}
}
