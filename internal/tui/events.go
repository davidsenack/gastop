package tui

import (
	"fmt"

	"github.com/davidsenack/gastop/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// EventsPanel displays the event log.
type EventsPanel struct {
	view     *tview.TextView
	maxLines int
	events   []model.Event
}

// NewEventsPanel creates a new events panel.
func NewEventsPanel(maxLines int) *EventsPanel {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(false)

	view.SetBorder(true).SetTitle(" EVENTS ")

	return &EventsPanel{
		view:     view,
		maxLines: maxLines,
	}
}

// Primitive returns the tview primitive.
func (p *EventsPanel) Primitive() tview.Primitive {
	return p.view
}

// Update updates the panel with new event data.
func (p *EventsPanel) Update(events []model.Event) {
	p.events = events

	var text string
	for _, e := range events {
		// Format: timestamp icon summary
		line := fmt.Sprintf("[gray]%s[-] %s %s",
			e.TimeString(),
			p.colorIcon(e.Icon(), e.Type),
			e.Summary(),
		)
		if e.Actor != "" {
			line += fmt.Sprintf(" [::d](%s)[-]", e.Actor)
		}
		text += line + "\n"
	}

	p.view.SetText(text)
	p.view.ScrollToEnd()
}

// colorIcon returns the icon with appropriate color.
func (p *EventsPanel) colorIcon(icon, eventType string) string {
	var color string
	switch eventType {
	case "spawn", "create", "bonded":
		color = "green"
	case "done", "completed", "merged":
		color = "green"
	case "crash", "failed", "merge_failed", "kill":
		color = "red"
	case "sling":
		color = "cyan"
	case "handoff":
		color = "yellow"
	case "nudge", "polecat_nudged":
		color = "yellow"
	case "in_progress":
		color = "blue"
	default:
		color = "white"
	}
	return fmt.Sprintf("[%s]%s[-]", color, icon)
}

// SetTitle sets the panel title.
func (p *EventsPanel) SetTitle(title string) {
	p.view.SetTitle(" " + title + " ")
}

// AppendEvent adds a new event to the display.
func (p *EventsPanel) AppendEvent(e model.Event) {
	p.events = append(p.events, e)
	if len(p.events) > p.maxLines {
		p.events = p.events[1:]
	}
	p.Update(p.events)
}

// SetBackgroundColor sets the background color.
func (p *EventsPanel) SetBackgroundColor(color tcell.Color) {
	p.view.SetBackgroundColor(color)
}
